package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	b, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.Write(b)
}

// GetEventHandler return the latest event and attendance info as json
func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	var b []byte
	w.Header().Set("content-type", "application/json")
	event, err := GetLatestEvent()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ = json.Marshal(map[string]string{
			"message": "could not find event",
		})
		w.Write(b)
		return
	}
	members, err := GetOKRsvpMembers(event.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ = json.Marshal(map[string]string{
			"message": "could not getting a member list",
		})
		w.Write(b)
		return
	}
	b, err = json.Marshal(members)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ = json.Marshal(map[string]string{
			"message": "could not getting a member list",
		})
		w.Write(b)
		return
	}
	w.Write(b)
}
