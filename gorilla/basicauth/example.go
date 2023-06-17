package basicauth

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Create a config with your desired values
	config := Config{
		Users: []User{
			{UserName: "user1", Password: "password1"},
		},
		RestrictedMethods: []string{"PUT", "POST", "PATCH", "DELETE", "GET"},
		RestrictedUrls:    []string{"/v1/user", "/v1/user/{key}", "/v1/admin"},
		RequireAuthForAll: true,
		UnauthorizedHandler: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		},
	}

	// Apply the Basic Auth middleware to the router
	router.Use(Middleware(config))

	// Define your routes and handlers
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the home page"))
	}).Methods("GET")

	router.HandleFunc("/v1/user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User endpoint"))
	}).Methods("GET")

	router.HandleFunc("/v1/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin endpoint"))
	}).Methods("GET")

	// Start the server
	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
