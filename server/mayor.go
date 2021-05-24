package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Mayor is the formated structure of a mayor
type Mayor struct {
	TermStart int    `json:"termStart"`
	TermEnd   int    `json:"termEnd"`
	Name      string `json:"name"`
	Current   bool   `json:"current"`
}

// getMayor returns a mayor for a given year. It is possible there was more
// than one mayor in a year, in which case it would provide you with both
// mayors at that time.
// GET /mayor/{year}
func getMayor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mayors []Mayor
	y := r.FormValue("year")
	year, e := strconv.Atoi(y)
	if e != nil {
		log.Println(fmt.Errorf("An error occurred while converting string to int: %v", e))
	}

	// out of bounds 1819-2020
	if year < 1819 || year > 2020 {
		fmt.Fprint(w, "Please provide a year between 1819 - 2020")
		return
	}

	for _, m := range Mayors {
		if m.TermStart <= year && m.Current {
			mayors = append(mayors, m)
			continue
		}
		if m.TermStart <= year && year <= m.TermEnd {
			mayors = append(mayors, m)
		}
	}

	json.NewEncoder(w).Encode(mayors)
}
