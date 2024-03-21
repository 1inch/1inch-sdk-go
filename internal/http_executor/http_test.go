package http_executor

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestExecuteRequest_Success(t *testing.T) {
	// Setup a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a dummy JSON object
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"result":"success"}`)
		if err != nil {
			t.Fatalf("Failed to writeString")
		}
	}))
	defer mockServer.Close()

	// Parse the mock server URL
	baseURL, _ := url.Parse(mockServer.URL)

	// Initialize the client with the mock server base URL
	client := &Client{
		baseURL:    baseURL,
		apiKey:     "testApiKey",
		httpClient: *mockServer.Client(),
	}

	// The data to be sent in the request
	data := RequestPayload{
		Method: "GET",
		Params: nil,
		U:      "/test",
		Body:   nil,
	}

	// The structure to store the response
	var result struct {
		Result string `json:"result"`
	}

	// Execute the request
	err := client.ExecuteRequest(context.Background(), data, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response body
	if result.Result != "success" {
		t.Errorf("Expected result 'success', got '%s'", result.Result)
	}
}

func TestExecuteRequest_SuccessfulPOST(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected 'POST', got '%s'", r.Method)
		}

		// Check for the presence of expected header
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected 'Content-Type' of 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Respond with a dummy JSON object
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status":"success"}`)
	}))
	defer mockServer.Close()

	client := Client{
		httpClient: *mockServer.Client(),
		baseURL:    mustParseURL(mockServer.URL), // Helper function to parse URL
		apiKey:     "testApiKey",
	}

	payload := RequestPayload{
		Method: "POST",
		U:      "/test",
		Body:   []byte(`{"key":"value"}`),
	}

	var response map[string]string
	if err := client.ExecuteRequest(context.Background(), payload, &response); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected response status 'success', got '%s'", response["status"])
	}
}

func TestExecuteRequest_ServerErrorPOST(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected 'POST', got '%s'", r.Method)
		}

		// Simulate server error
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"message":"internal server error"}`)
	}))
	defer mockServer.Close()

	client := Client{
		httpClient: *mockServer.Client(),
		baseURL:    mustParseURL(mockServer.URL), // Helper function to parse URL
		apiKey:     "testApiKey",
	}

	payload := RequestPayload{
		Method: "POST",
		U:      "/error", // Endpoint to trigger an error
		Body:   []byte(`{"key":"value"}`),
	}

	var response map[string]interface{} // Using interface{} to potentially capture any structure of response
	err := client.ExecuteRequest(context.Background(), payload, &response)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	// Assuming your error handling includes parsing the error response and including it in the error message
	expectedErrorMessage := "processing response failed: HTTP 500: internal server error"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestAuthorizationKey(t *testing.T) {
	// Expected API key value
	expectedAPIKey := "testApiKey"

	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+expectedAPIKey {
			t.Errorf("Authorization header does not match expected. Got %s, want %s", authHeader, "Bearer "+expectedAPIKey)
		}

		// If the Authorization header is as expected, respond with 200 OK
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status":"success"}`)
	}))
	defer mockServer.Close()

	client := Client{
		httpClient: *mockServer.Client(),
		baseURL:    mustParseURL(mockServer.URL),
		apiKey:     expectedAPIKey,
	}

	payload := RequestPayload{
		Method: "GET", // Method can be anything for this test
		U:      "/",   // Endpoint can be anything for this test
	}

	var response interface{}

	err := client.ExecuteRequest(context.Background(), payload, &response)
	if err != nil {
		t.Fatalf("Did not expect an error, got %v", err)
	}
}

// Helper function to parse URL and ensure no error is returned
func mustParseURL(u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url
}
