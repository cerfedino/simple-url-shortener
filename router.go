package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	var model Storage = &postgresstorage{}
	model.init("")

	r := mux.NewRouter()

	// Log all requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Redirect '/' to '/static/index.html'
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	// Serve static files
	r.Methods("GET").PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Methods("GET").Path("/{shortURL}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortURL := vars["shortURL"]
		longURL, err := model.getLongerURL(shortURL)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/static/error.html?e=404", http.StatusFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	})

	// Catch every other request
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/error.html?e=404", http.StatusFound)
	})

	return r
}
