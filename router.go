package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type conf struct {
	Trusted_proxies []string `yaml:"trusted_proxies"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Printf("%+v\n", c)
	return c
}

func ContainsString(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}
	return false
}

func Router() *mux.Router {
	var model Storage = &postgresstorage{}
	model.init("")

	var cc conf
	cc.getConf()

	publicFileServeHandler := http.StripPrefix("/public/", http.FileServer(http.Dir("./static/public")))
	privateFileServeHandler := http.StripPrefix("/private/", http.FileServer(http.Dir("./static/private")))

	// model.removeShortenedURL("private")
	// model.storeShortenedURL("/private/secret.html", false, "private")
	// fmt.Println(model.getAllShortenedURLs(false))

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			real_ip := r.RemoteAddr
			// Strip the port from the remote address
			if strings.Count(real_ip, ":") == 1 || strings.Count(real_ip, ":") == 8 {
				real_ip = real_ip[:strings.LastIndex(real_ip, ":")]
			}

			xff := r.Header.Get("X-Forwarded-For")
			// log.Println(fmt.Sprintf("X-Forwarded-For: %s", xff))
			if xff != "" && ContainsString(cc.Trusted_proxies, real_ip) {
				xffs := strings.Split(xff, ",")
			out:
				for i := len(xffs) - 1; i >= 0; i-- {
					xffs[i] = strings.TrimSpace(xffs[i])
					if !ContainsString(cc.Trusted_proxies, xffs[i]) {
						real_ip = xffs[i]
						break out
					}
				}
			}
			// log.Println(fmt.Sprintf("real_ip: %s", real_ip))

			// Save the real IP address in the request context
			r = r.WithContext(context.WithValue(r.Context(), "real_ip", real_ip))
			next.ServeHTTP(w, r)
		})
	})

	// Log all requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Context().Value("real_ip").(string), r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Matcher for both '/' and '/index.html'
	r.Methods("GET").Path("/{path:index.html|}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/public/index.html", http.StatusFound)
	})

	// Serve static files in the static/public directory
	r.Methods("GET").PathPrefix("/public/").Handler(publicFileServeHandler)

	// Attempt to redirect to a long url
	r.Methods("GET").Path("/{shortURL}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortURL := vars["shortURL"]
		// Retrieve the long URL from the database
		longURL, err := model.getLongerURL(shortURL)
		defer model.logShorteningRequest(r.Context().Value("real_ip").(string), shortURL, longURL)
		// If the short URL is not found, redirect to the error page
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, fmt.Sprintf("/public/error.html?e=%d&reqURL=%v", http.StatusNotFound, shortURL), http.StatusFound)
			return
		}
		// Instead if the URL redirects to the static/private directory, serve the file
		if strings.HasPrefix(longURL, "/private/") {
			r.URL.Path = longURL
			privateFileServeHandler.ServeHTTP(w, r)
			return
		}
		// Otherwise, redirect to the long URL
		http.Redirect(w, r, longURL, http.StatusFound)
	})

	// Catch every other request
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/public/error.html?e=%d&reqURL=%v", http.StatusNotFound, r.URL), http.StatusFound)
	})

	return r
}
