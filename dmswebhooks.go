// Package dmswebhooks provides an http.Handler that makes it easy to receive
// webhooks from Dead Man's Snitch's webhook integration. This makes building
// custom integrations easier.
package dmswebhooks

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	TypeSnitchReporting = "snitch.reporting"
	TypeSnitchErrored   = "snitch.errored"
	TypeSnitchMissing   = "snitch.missing"

	StatusPending = "pending"
	StatusHealthy = "healthy"
	StatusMissing = "missing"
	StatusErrored = "errored"
)

// Alert provides information about updates to Snitch state.
type Alert struct {
	// Type gives information about what kind of alert occurred.
	Type string `json:"type"`

	// Timestamp is when the alert occurred.
	Timestamp time.Time `json:"timestamp"`

	// Data contains information about the Snitch
	Data struct {
		Snitch struct {
			Token           string   `json:"token"`
			Name            string   `json:"name"`
			Notes           string   `json:"notes"`
			Tags            []string `json:"tags"`
			status          string   `json:"status"`
			previous_status string   `json:"previous_status"`
		} `json:"snitch"`
	} `json:"data"`
}

type handler struct {
	callback func(*Alert) error
}

// NewHandler returns an http.Handler that parses incoming webhooks from Dead
// Man's Snitch's webhook integration and calls the given callback handler.
// Errors returned from the callback are sent back to the calling server
// causing the request to be retried.
func NewHandler(callback func(*Alert) error) http.Handler {
	return &handler{
		callback: callback,
	}
}

func (h *handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	alert := &Alert{}

	// Dead Man's Snitch will always use POST for webhooks
	if req.Method != "POST" {
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(alert); err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		// This error is internal to the library so it should be safe to send back
		// without exposing any internal data. It may mean we can correct it server
		// side as well.
		res.Write([]byte(err.Error()))

		return
	}

	if err := h.callback(alert); err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		// Avoid sending client error information back to the server to avoid
		// possibly exposing anything sensitive.
		res.Write([]byte("Error in callback"))

		return
	}

	// Yeah! Everything worked!
	res.WriteHeader(http.StatusNoContent)
}
