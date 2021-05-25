package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Mayor is the formated structure of a mayor
type Mayor struct {
	TermStart int    `json:"termStart"`
	TermEnd   int    `json:"termEnd"`
	Name      string `json:"name"`
	Current   bool   `json:"current"`
}

type yearCtxKey struct{}

func getCurrentMayor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	for _, m := range Mayors {
		if m.Current {
			if err := json.NewEncoder(w).Encode(m); err != nil {
				log.Println("Error while encoding current mayor:", err)
			}
			break
		}
	}
}

// getMayor returns a mayor for a given year. It is possible there was more
// than one mayor in a year, in which case it would provide you with both
// mayors at that time.
// GET /mayor/{year}
func getMayor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mayors []*Mayor
	y := r.Context().Value(yearCtxKey{}).(int)

	for _, m := range Mayors {
		if m.TermStart <= y && m.Current {
			mayors = append(mayors, m)
			continue
		}
		if m.TermStart <= y && y <= m.TermEnd {
			mayors = append(mayors, m)
		}
	}

	if err := json.NewEncoder(w).Encode(mayors); err != nil {
		log.Println("Error while encoding current mayor:", err)
	}
}

// MayorCtx is custom middleware to load a Mayor object, from the url
func MayorCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		y := chi.URLParam(r, "year")
		year, err := strconv.Atoi(y)
		if err != nil {
			log.Println(fmt.Errorf("An error occurred while converting string to int: %v", err))
			err := render.Render(w, r, &ErrResponse{
				StatusText:     "Invalid year provided",
				HTTPStatusCode: 400,
				ErrorText:      fmt.Sprintf("An error occurred while converting string to int: %v", err),
			})
			if err != nil {
				fmt.Println("Error rendering custom error object:", err)
			}
			return
		}
		// out of bounds 1819-2020
		if year < 1819 || year > 2020 {
			err := render.Render(w, r, &ErrResponse{
				StatusText:     "Please provide a year between 1819 - 2020",
				HTTPStatusCode: 400,
				ErrorText:      fmt.Sprintf("Invalid year was passed in: %v", year),
			})
			if err != nil {
				fmt.Println("Error creating custom error object:", err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), yearCtxKey{}, year)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
