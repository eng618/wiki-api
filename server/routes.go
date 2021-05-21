package server

import (
	"fmt"
	"net/http"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from wiki server")
}

func errorTest(w http.ResponseWriter, r *http.Request) {
	panic("OH NO 😵")
}
