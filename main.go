package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const port = "8000"

func main() {
	http.HandleFunc("/", IndexView)
	http.HandleFunc("/event", GetEventHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Listening on %s.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// IndexView render the index template
func IndexView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	f, err := os.Open("index.html")
	if err != nil {
		panic(err)
	}
	io.Copy(w, f)
}

// GetEventHandler return the latest event and attendance info as json
func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	event, err := GetLatestEvent()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}
	members, err := GetOKRsvpMembers(event.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}
	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}
}
