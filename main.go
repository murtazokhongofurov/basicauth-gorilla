package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// This is configuration struct of Basic Auth
type Config struct {
	// Users is list of users that have access.
	// There may be a user1 with password1 and user2 with password2 and etc.
	Users []User `json:"users"`
	// Restricted Method means that the middleware only applies for method are given.
	// For example, PUT, POST, PATCH, DELETE methods are given to this field. Middleware check password for request with these REST Methods.
	RestrictedMethods []string `json:"restricted_methods"`
	// Restricted urls are the urls that are authoriztion is required.
	// For example, /v1/user, /v1/user/{key}, /v1/admin
	// if /v1/user is given, request url is checked for equality.
	// if /v1/user/{key} is given, request url is check for tthe urls starting with /v1/user  and one other key.
	// if /v1/user/*  is given, request url is check for all the urls that starts with '/v1/user'.
	RestrictedUrls []string `json:"restricted_urls"`
	// If this field is set to true, all the requests are authenticated
	// If this field is not set or set to true, other fields are checked such as, RestrictedMethods and RestrictedUrls
	RequireAuthForAll   bool `json:"require_auth_for_all"`
	
	UnauthorizedHandler http.HandlerFunc
}

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func New(config Config) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !requiresAuth(config, r) {
				h.ServeHTTP(w, r)
				return
			}
			username, password, ok := r.BasicAuth()
			if !ok || !isAuthorized(username, password, config.Users) {
				config.UnauthorizedHandler(w, r)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func requiresAuth(config Config, r *http.Request) bool {
	if config.RequireAuthForAll {
		return true
	}

	method := r.Method
	for _, restrictedMethod := range config.RestrictedMethods {
		if method == restrictedMethod {
			return true
		}
	}

	url := r.URL.Path
	for _, restrictedURL := range config.RestrictedUrls {
		if matchURL(restrictedURL, url) {
			return true
		}
	}

	return false
}

func matchURL(pattern, url string) bool {
	if pattern == url {
		return true
	}

	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(url, prefix)
	}

	return false
}

func isAuthorized(username, password string, users []User) bool {
	for _, user := range users {
		if user.UserName == username && user.Password == password {
			return true
		}
	}
	return false
}
