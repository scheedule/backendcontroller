package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/scheedule/backend_controller/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Server struct {
	session_name   string
	session_secret string
	session_store  *sessions.CookieStore
	services       map[string]*url.URL
}

func New(sessionName, sessionSecret string, services map[string]*url.URL) *Server {

	srv := &Server{
		session_name:   sessionName,
		session_secret: sessionSecret,
		session_store:  sessions.NewCookieStore([]byte(sessionSecret)),
		services:       services,
	}

	// Proxy services
	prx := proxy.New(srv.services, srv.isAuth)
	http.Handle("/prx/", http.StripPrefix("/prx", prx))

	// Login handle
	http.HandleFunc("/oauth/", srv.oauthCallback)

	// Static GUI
	http.Handle("/", http.FileServer(http.Dir("./public")))

	return srv
}

// Callback to interrogate OAuth token in session
func (s *Server) oauthCallback(w http.ResponseWriter, r *http.Request) {
	session, err := s.session_store.Get(r, s.session_name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.FormValue("token")
	log.Debug("Received token", token)

	if token == "" {
		http.NotFound(w, r)
	} else { // New Token!
		fmt.Println("New token")
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
		fmt.Println("Done..:", m["sub"])
	}
}

func (s *Server) isAuth(r *http.Request) bool {
	fmt.Println("Checking auth")

	session, err := s.session_store.Get(r, s.session_name)
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