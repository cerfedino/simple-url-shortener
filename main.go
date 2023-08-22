package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var model Storage = &postgresstorage{}
	model.init("")

	// model.storeShortenedURL("https://www.google.com", "tg")

	// long_url, err := model.getLongerURL("tg")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(long_url)

	// model.removeShortenedURL("tg")
	// _, err = model.getLongerURL("tg")
	// if err != nil {
	// 	log.Println(err)
	// }

	srv := &http.Server{
		Handler: Router(),
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		} else {
			log.Printf("Listening on %s ...\n", srv.Addr)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5000)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down ...")
	os.Exit(0)
}
