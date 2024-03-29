package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Quantaly/gdqIcal/lib"
)

func main() {
	_, present := os.LookupEnv("DYNO") // detect running on Heroku
	if present {
		log.SetFlags(0)
	}

	port, present := os.LookupEnv("PORT")
	if !present {
		log.Fatal("PORT environment variable not set")
	}

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

	http.HandleFunc("/yeet", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("yeet\n"))
		if err != nil {
			log.Println("failed to yeet: ", err)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
