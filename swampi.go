package swampi

import (
	"io"
	"net/http"
)

// Swampi is the swarm http api client
type Swampi struct {
	c *http.Client
	// the bare host (http://foobar:8500)
	endpoint string
}

// New creates a new instance of swampi
func New(endpoint string) *Swampi {
	return &Swampi{
		c:        &http.Client{},
		endpoint: endpoint,
	}
}

// Do executes the given http request
func (s *Swampi) Do(req *http.Request) (*http.Response, error) {
	return s.c.Do(req)
}

// Send is a wrapper around constructCall + Do
func (s *Swampi) Send(call APICall, body io.Reader, headers map[string][]string, args ...interface{}) (*http.Response, error) {
	req, err := s.constructCall(call, body, headers, s.endpoint+call.ParseArgs(args...))
	if err != nil {
		return nil, err
	}
	return s.Do(req)
}

func (s *Swampi) constructCall(call APICall, body io.Reader, headers map[string][]string, url string) (*http.Request, error) {
	req, err := http.NewRequest(call.Method(), url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		if len(v) == 1 {
			req.Header.Set(k, v[0])
		} else {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
	}
	return req, nil
}
