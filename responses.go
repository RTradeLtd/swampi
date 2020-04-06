package swampi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"net/http"
)

// contains structs for complex responses returned by the swarm api

// Response is a struct that bundles a http response
// and some helper functinos
type Response struct {
	call APICall
	*http.Response
}

// BZZList is a struct containing responses from the ListFiles API
type BZZList struct {
	Entries []struct {
		Hash        string    `json:"hash"`
		Path        string    `json:"path"`
		ContentType string    `json:"contentType"`
		Mode        int       `json:"mode"`
		Size        int       `json:"size"`
		ModTime     time.Time `json:"mod_time"`
	} `json:"entries"`
}

// SwarmUnmarshal is a helper method to unfurl the response returned from the api
func (r *Response) SwarmUnmarshal() (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	out := r.call.Response()
	// if out is nil this indicates the response from the call is a simple, single value
	// for example the upload call returns just the hash, and theres no need to complicate it
	if out == nil {
		return nil, errors.New("api call has no specific response associated with it")
	}
	if err := json.Unmarshal(data, out); err != nil {
		return nil, err
	}
	return out, nil
}
