package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	var model Storage = &postgresstorage{}
	model.init("")

	publicFileServeHandler := http.StripPrefix("/public/", http.FileServer(http.Dir("./static/public")))
	privateFileServeHandler := http.StripPrefix("/private/", http.FileServer(http.Dir("./static/private")))

	// model.removeShortenedURL("private")
	// model.storeShortenedURL("/private/secret.html", false, "private")
	// fmt.Println(model.getAllShortenedURLs(false))

	r := mux.NewRouter()

	// Log all requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Matcher for both '/' and '/index.html'
	r.Methods("GET").Path("/{path:index.html|}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/public/index.html", http.StatusFound)
	})

	// Serve static files
	r.Methods("GET").PathPrefix("/public/").Handler(publicFileServeHandler)

	// Shorten URL
	r.Methods("GET").Path("/{shortURL}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortURL := vars["shortURL"]
		longURL, err := model.getLongerURL(shortURL)
		defer model.logShorteningRequest(r.RemoteAddr, shortURL, longURL)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, fmt.Sprintf("/public/error.html?e=%d&reqURL=%v", http.StatusNotFound, shortURL), http.StatusFound)
			return
		}
		// If the URL redirects to the private directory, serve the file
		if strings.HasPrefix(longURL, "/private/") {
			r.URL.Path = longURL
			privateFileServeHandler.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	})

	// Catch every other request
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/public/error.html?e=%d&reqURL=%v", http.StatusNotFound, r.URL), http.StatusFound)
	})

	return r
}
