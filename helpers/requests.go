package helpers

import (
	"bytes"
	"net/http"
)

func CreateAPIRequest(method string, url string, headers map[string]string, body string) *http.Request {
	// Create a new Request-Element with all top parameters
	request, err := http.NewRequest(method, url, bytes.NewBufferString(body))

	Fatal("[helpers].CreateAPIRequest:1", err)

	// Set headers and return the Request-Object
	request.Header["Content-Type"] = []string{"application/json"}

	for key, val := range headers {
		request.Header[key] = []string{val}
	}

	return request
}
