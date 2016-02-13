// Package server provides all handlers and session management for the application.
package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/scheedule/backendcontroller/proxy"
)

// Type Server holds session information and also which services are registered
// to be proxied to.
type Server struct {
	sessionName   string
	sessionSecret string
	sessionStore  *sessions.CookieStore
	services      map[string]*url.URL
	public        bool
}

// Create a new server with session configuration and service registry
func New(sessionName, sessionSecret string, services map[string]*url.URL, public bool) *Server {

	srv := &Server{
		sessionName:   sessionName,
		sessionSecret: sessionSecret,
		sessionStore:  sessions.NewCookieStore([]byte(sessionSecret)),
		services:      services,
		public:        public,
	}

	// Proxy services
	prx := proxy.New(srv.services, srv.isAuth)
	http.Handle("/prx/", http.StripPrefix("/prx", prx))

	// Login handle
	http.HandleFunc("/oauth/", srv.oauthCallback)

	// Static GUI
	http.Handle("/", http.FileServer(http.Dir("./public")))

	log.Info("Handling web GUI on /")
	log.Info("Handling proxying to services on /prx/")

	return srv
}

// Callback to interrogate OAuth token in session
func (s *Server) oauthCallback(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, s.sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.FormValue("token")
	log.Debug("Received token", token)

	if token == "" {
		http.NotFound(w, r)
	} else { // New Token!
		log.Debug("New token")
		res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var m map[string]string
		json.Unmarshal(body, &m)
		session.Values["user_id"] = m["sub"]
		session.Values["name"] = m["name"]
		session.Save(r, w)
		log.Debug("Done..:", m["sub"])
		w.WriteHeader(200)
	}
}

// Authentication method to be used by proxy. Checks session for user_id variable
// set.
func (s *Server) isAuth(r *http.Request) bool {
	log.Debug("Checking auth")
	log.Debug("Public: ", s.public)

	if s.public {
		r.Header.Set("user_id", "")
		return true
	}

	session, err := s.sessionStore.Get(r, s.sessionName)
	if err != nil {
		return false
	}

	// Get session value
	userIdInterface := session.Values["user_id"]

	// Set header with user_id
	if userId, ok := userIdInterface.(string); ok {
		r.Header.Set("user_id", userId)
		return true
	}
	return false
}
