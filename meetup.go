package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// MeetupError represents an error of meetup API
type MeetupError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e MeetupError) Error() string {
	return e.Message
}

// MeetupErrors contains a list of MeetupErrors
type MeetupErrors struct {
	Errors []MeetupError `json:"errors"`
}

func (e MeetupErrors) Error() string {
	errs := make([]string, 0)
	for _, err := range e.Errors {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, " + ")
}

// MeetupEvent represents an event with minimal info
type MeetupEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LocalDate   string `json:"local_date"`
	LocalTime   string `json:"local_time"`
	OKRsvpCount int    `json:"yes_rsvp_count"`
}

// MeetupMember represents a member with minimal info
type MeetupMember struct {
	Member struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Photo struct {
			ID        int    `json:"id"`
			PhotoLink string `json:"photo_link"`
		} `json:"photo"`
	} `json:"member"`
	Response string `json:"response"`
}

const (
	endpoint = "http://api.meetup.com/GDG-Golang-Korea"
)

func meetupBuildURL(path string) string {
	return endpoint + path
}

func meetupRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetLatestEvent get a latest event information
func GetLatestEvent() (*MeetupEvent, error) {
	url := meetupBuildURL("/events?status=past,upcoming&desc=true&page=1")
	b, err := meetupRequest(url)
	if err != nil {
		return nil, err
	}
	var events []MeetupEvent
	err = json.Unmarshal(b, &events)
	if err := json.Unmarshal(b, &events); err != nil {
		var aerr MeetupErrors
		if err := json.Unmarshal(b, &aerr); err == nil {
			return nil, aerr
		}
		return nil, err
	}
	return &events[0], nil
}

// GetOKRsvpMembers get a list of confirmed members of an event
func GetOKRsvpMembers(id string) ([]MeetupMember, error) {
	url := meetupBuildURL(fmt.Sprintf("/events/%s/rsvps", id))
	b, err := meetupRequest(url)
	if err != nil {
		return nil, err
	}
	var members []MeetupMember
	if err := json.Unmarshal(b, &members); err != nil {
		var aerr MeetupErrors
		if err := json.Unmarshal(b, &aerr); err == nil {
			return nil, aerr
		}
		return nil, err
	}
	var okMembers []MeetupMember
	for _, m := range members {
		if m.Response == "yes" {
			okMembers = append(okMembers, m)
		}
	}
	return okMembers, nil
}
