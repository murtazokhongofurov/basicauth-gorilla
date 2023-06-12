package basicauth

import (
	"net/http"
	"strings"
)

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
