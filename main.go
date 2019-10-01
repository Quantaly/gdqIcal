package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/Quantaly/gdqIcal/lib"
)

func main() {
	log.Println("Preparing to serve")

	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	http.HandleFunc("/out.ics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/calendar")
		w.WriteHeader(200)
		lib.GenerateCalendar(w, "https://gamesdonequick.com/schedule/27") // GDQx 2019
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
