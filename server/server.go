// Package server manages the wiki server
package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Server is a struct that contains, all the required pieces to effectively run
// the main server and handles all the routing to handle different routes.
type Server struct {
	Router  *mux.Router
	Negroni *negroni.Negroni
}

func (s *Server) init() {
	// Initialize Routes
	s.Router = mux.NewRouter().StrictSlash(true)
	s.initializeRoutes()

	// Initialize Middleware (logging)
	s.Negroni = negroni.Classic()
	s.Negroni.UseHandler(s.Router)
}

// Run starts the API server
func (s *Server) Run(addr string) {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0%s", addr),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.Router,
	}

	// launch server as go routine so that we can continue the control flow
	go func() {
		log.Printf("Server is listening at : 127.0.0.1%v\n", addr)
		if err := http.ListenAndServe(addr, s.Negroni); err != nil {
			log.Println(err)
		}
	}()

	// set up channel to block until we have shutdown signal from os
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	go func() {
		log.Println("initiating shutdown, please wait...")
		srv.Shutdown(ctx)
	}()
	<-ctx.Done()

	log.Println("shutting down")
	os.Exit(0)
}

func (s *Server) initializeRoutes() {
	// home
	s.Router.HandleFunc("/", homeLink)
	s.Router.HandleFunc("/panic", errorTest)
}
