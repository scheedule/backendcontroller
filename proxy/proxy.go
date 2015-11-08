package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewUnauthorizedResponse(r *http.Request) *http.Response {
	resp := &http.Response{}
	resp.Header = make(http.Header)
	resp.StatusCode = http.StatusUnauthorized
	buf := bytes.NewBufferString("")
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

type AuthorizedTransport struct {
	RoundTripper http.RoundTripper
	IsAuth       func(*http.Request) bool
}

func (at *AuthorizedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if at.IsAuth(req) {
		fmt.Println("Letting round trip")
		return at.RoundTripper.RoundTrip(req) // Let them through
	} else {
		fmt.Println("Returning unauthorized")
		return NewUnauthorizedResponse(req), nil
	}
}

func New(targets map[string]*url.URL, isAuth func(*http.Request) bool) *httputil.ReverseProxy {
	director := func(r *http.Request) {
		fmt.Println(r.URL.Path)

		spl := strings.Split(r.URL.Path, "/")
		fmt.Println(spl)

		if len(spl) < 2 {
			fmt.Println("Failed to match")
			return
		}

		service := spl[1]
		fmt.Println("Service:", service)

		target := targets[service]
		if target == nil {
			fmt.Println("Failed to lookup service")
			return
		}
		fmt.Println("url:", target)

		newPath := "/" + strings.Join(spl[2:], "/")

		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.URL.Path = newPath
	}

	transport := &AuthorizedTransport{http.DefaultTransport, isAuth}

	return &httputil.ReverseProxy{Director: director, Transport: transport}
}
