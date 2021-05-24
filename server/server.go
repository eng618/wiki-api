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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Server is a struct that contains, all the required pieces to effectively run
// the main server and handles all the routing to handle different routes.
type Server struct {
	Router *chi.Mux
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render formats the custom error object with a http status code
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrNotFound is a standard error for not found.
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

// Initialize prepares a server along with all middleware.
func (s *Server) Initialize() {
	// Initialize Routes
	s.Router = chi.NewRouter()
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.URLFormat)
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))

	s.initializeRoutes()
}

// Run starts the API server
func (s *Server) Run(addr string) {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*5, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
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
		if err := http.ListenAndServe(addr, s.Router); err != nil {
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
	r := s.Router
	// home
	r.Get("/", homeLink)
	r.Get("/panic", errorTest)

	r.Route("/mayor", func(r chi.Router) {
		r.Get("/", getCurrentMayor) // GET /mayor
		r.Route("/{year}", func(r chi.Router) {
			r.Use(MayorCtx)
			r.Get("/", getMayor)
		})
	})
}
