package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

// setup sets up a test HTTP server along with a Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (*Client, *http.ServeMux, string, func(), error) {
	apiHandler := http.NewServeMux()

	// This defaults all requests to return a 404
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
	})

	// server is a test HTTP server used to provide mock API responses.
	// the base URL of the client will have its destination swapped to use this new test server for requests
	server := httptest.NewServer(apiHandler)
	c, err := NewClient(
		Config{
			TargetEnvironment: EnvironmentProduction,
			ApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
		})
	if err != nil {
		return nil, nil, "", nil, err
	}

	url, _ := url.Parse(server.URL + "/")
	c.BaseURL = url
	return c, apiHandler, server.URL, server.Close, nil
}
