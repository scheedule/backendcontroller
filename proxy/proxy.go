// Package proxy provides a reverse proxy to several backend services. The proxy
// allows configuration of authentication filtering.
package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Create HTTP response with status unauthorized
func NewUnauthorizedResponse(r *http.Request) *http.Response {
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.StatusCode = http.StatusUnauthorized
	buf := bytes.NewBufferString("")
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

// Transport configurable with authentication
type AuthorizedTransport struct {
	RoundTripper http.RoundTripper
	IsAuth       func(*http.Request) bool
}

// Transport interface method to process request and return response. IF we the
// request fails to authenticate, a response indicating status unauthorized
// is returned. If the request is authenticated, the request continues on to
// the proper service.
func (at *AuthorizedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if at.IsAuth(req) {
		log.Debug("letting round trip")
		return at.RoundTripper.RoundTrip(req) // Let them through
	} else {
		log.Debug("returning unauthorized")
		return NewUnauthorizedResponse(req), nil
	}
}

// Create new proxy. Proxy can be configured with "isAuth" argument.
func New(targets map[string]*url.URL, isAuth func(*http.Request) bool) *httputil.ReverseProxy {
	director := func(r *http.Request) {
		log.Debug(r.URL.Path)

		spl := strings.Split(r.URL.Path, "/")

		if len(spl) < 2 {
			log.Warn("failed to match path: ", spl)
			return
		}

		service := spl[1]
		log.Debug("directing to service: ", service)

		target, ok := targets[service]
		if !ok {
			log.Warn("failed to lookup service")
			return
		}

		newPath := "/" + strings.Join(spl[2:], "/")

		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.URL.Path = newPath
	}

	transport := &AuthorizedTransport{http.DefaultTransport, isAuth}

	return &httputil.ReverseProxy{Director: director, Transport: transport}
}
