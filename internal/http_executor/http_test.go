package http_executor

import (
	"bytes"
	"context"
	"encoding/json"
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
		Params: "/test",
		Body:   nil,
	}

	// The structure to store the response
	var result struct {
		Result string `json:"result"`
	}

	// Execute the request
	resp, err := client.ExecuteRequest(context.Background(), data, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	if result.Result != "success" {
		t.Errorf("Expected result 'success', got '%s'", result.Result)
	}
}

func TestExecuteRequest_PostWithBody(t *testing.T) {
	// Define the request body
	requestBody := map[string]string{"key": "value"}
	encodedRequestBody, _ := json.Marshal(requestBody)

	// Setup a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Verify the request body
		bodyBytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		if !bytes.Equal(bodyBytes, encodedRequestBody) {
			t.Errorf("Expected request body '%s', got '%s'", string(encodedRequestBody), string(bodyBytes))
		}

		// Respond with a dummy JSON object
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"response":"processed"}`)
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
		Method: "POST",
		Params: "/submit",
		Body:   encodedRequestBody,
	}

	// The structure to store the response
	var response struct {
		Response string `json:"response"`
	}

	// Execute the request
	resp, err := client.ExecuteRequest(context.Background(), data, &response)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	if response.Response != "processed" {
		t.Errorf("Expected response 'processed', got '%s'", response.Response)
	}
}
