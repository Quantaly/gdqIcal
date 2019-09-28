package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/Quantaly/gdqIcal/lib"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	log.Println("Preparing to serve")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", struct{}{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/out.ics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/calendar")
		w.WriteHeader(200)
		lib.GenerateCalendar(w, "https://gamesdonequick.com/schedule/27") // GDQx 2019
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
