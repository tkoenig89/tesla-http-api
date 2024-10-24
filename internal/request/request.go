package request

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

type Endpoint struct {
	Url                string
	Method             string
	Headers            map[string]string
	Body               string
	InsecureSkipVerify bool
}

func Connect(e Endpoint) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: e.InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(e.Method, e.Url, strings.NewReader(e.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to endpoint: %s", err)
	}

	for k, v := range e.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %s", err)
	}
	return resp, nil
}