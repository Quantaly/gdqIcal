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

	scheduleURL, present := os.LookupEnv("SCHEDULE_URL")
	if !present {
		scheduleURL = "https://gamesdonequick.com/schedule"
	}

	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	http.HandleFunc("/out.ics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/calendar")
		w.WriteHeader(200)
		lib.GenerateCalendar(w, scheduleURL)
		go lib.LogRequest(r.Header.Get("X-Forwarded-For"), "out.ics")
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
