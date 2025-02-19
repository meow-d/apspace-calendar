package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/meow-d/apspace-calendar/src/calendar"
)

func handler(w http.ResponseWriter, r *http.Request) {
	intake := r.URL.Query().Get("intake")
	if intake == "" {
		http.Error(w, "Missing required parameter: intake", http.StatusBadRequest)
		return
	}

	titleFormat := r.URL.Query().Get("title")
	if titleFormat == "" {
		titleFormat = "module_name"
	}

	icsData, err := calendar.FetchAndConvert(intake, titleFormat)
	if err != nil {
		http.Error(w, "Failed to process calendar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
	w.Write([]byte(icsData))
}

func main() {
	serve := flag.Bool("serve", false, "Run server")
	flag.Parse()

	if *serve {
		http.HandleFunc("/", handler)
		fmt.Println("Server running at http://localhost:8080")
		http.ListenAndServe(":8080", nil)
	}
}
