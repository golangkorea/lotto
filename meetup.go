package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// APIError represents an error of meetup API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

// APIErrors contains a list of APIErrors
type APIErrors struct {
	Errors []APIError `json:"errors"`
}

func (e APIErrors) Error() string {
	errs := make([]string, 0)
	for _, err := range e.Errors {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, " + ")
}

// SimpleEvent represents an event with minimal info
type SimpleEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LocalDate   string `json:"local_date"`
	LocalTime   string `json:"local_time"`
	OKRsvpCount int    `json:"yes_rsvp_count"`
}

// SimpleMember represents a member with minimal info
type SimpleMember struct {
	Member struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Photo struct {
			ID        int    `json:"id"`
			PhotoLink string `json:"photo_link"`
		} `json:"photo"`
	} `json:"member"`
	RSVP struct {
		ID       int    `json:"id"`
		Response string `json:"response"`
	} `json:"rsvp"`
}

const (
	endpoint = "http://api.meetup.com/GDG-Golang-Korea"
)

func buildURL(path string) string {
	return endpoint + path
}

func request(url string) ([]byte, error) {
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
func GetLatestEvent() (*SimpleEvent, error) {
	url := buildURL("/events?status=past,upcoming&desc=true&page=1")
	b, err := request(url)
	if err != nil {
		return nil, err
	}
	var events []SimpleEvent
	err = json.Unmarshal(b, &events)
	if err := json.Unmarshal(b, &events); err != nil {
		var aerr APIErrors
		if err := json.Unmarshal(b, &aerr); err == nil {
			return nil, aerr
		}
		return nil, err
	}
	return &events[0], nil
}

// GetOKRsvpMembers get a list of confirmed members of an event
func GetOKRsvpMembers(id string) ([]SimpleMember, error) {
	url := buildURL(fmt.Sprintf("/events/%s/attendance", id))
	b, err := request(url)
	if err != nil {
		return nil, err
	}
	var members []SimpleMember
	if err := json.Unmarshal(b, &members); err != nil {
		var aerr APIErrors
		if err := json.Unmarshal(b, &aerr); err == nil {
			return nil, aerr
		}
		return nil, err
	}
	var okMembers []SimpleMember
	for _, m := range members {
		if m.RSVP.Response == "yes" {
			okMembers = append(okMembers, m)
		}
	}
	return okMembers, nil
}
