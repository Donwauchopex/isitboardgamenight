package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	cancelled bool
	auth      string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	v, ok := os.LookupEnv("AUTHORIZATION")
	if !ok {
		panic("AUTHORIZATION environment variable not set")
	}
	auth = v
}

func nextBoardGameNight(t time.Time) time.Time {
	boardGameNight := t

	for boardGameNight.Weekday() != time.Tuesday {
		boardGameNight = boardGameNight.AddDate(0, 0, 1)
	}

	boardGameNight = time.Date(
		boardGameNight.Year(),
		boardGameNight.Month(),
		boardGameNight.Day(),
		18,
		45,
		0,
		0,
		boardGameNight.Location(),
	)

	return boardGameNight
}

func index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if cancelled {
		w.Write([]byte("Board game night has been cancelled :("))
		return
	}

	now := time.Now()
	nextBoardGameNight := nextBoardGameNight(now)

	if now.After(nextBoardGameNight) &&
		now.Before(nextBoardGameNight.Add(5*time.Hour+14*time.Minute)) {
		w.Write([]byte("It is board game night!"))
		return
	}

	diff := nextBoardGameNight.Sub(now)

	w.Write(
		[]byte(
			fmt.Sprintf(
				"The next board game night is on %s at 6:45PM\n",
				nextBoardGameNight.Format("2024-01-02"),
			),
		),
	)
	w.Write(
		[]byte(
			fmt.Sprintf(
				"That is in %d days, %d hours, %d minutes, and %d seconds\n",
				int(diff.Hours()/24),
				int(diff.Hours())%24,
				int(diff.Minutes())%60,
				int(diff.Seconds())%60,
			),
		),
	)
}

func updateBoardGameNight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Header.Get("Authorization") != auth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized."))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Invalid request method."))
	}

	cancelled = !cancelled
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Board game night has been updated!"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/update", updateBoardGameNight)
	http.ListenAndServe(":8080", nil)
}
