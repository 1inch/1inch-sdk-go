package http_executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
)

func TestExecuteRequest_SuccessGET(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"result":"success"}`)
		if err != nil {
			t.Fatalf("Failed to writeString")
		}
	}))
	defer mockServer.Close()

	baseURL, _ := url.Parse(mockServer.URL)

	client := &Client{
		baseURL:    baseURL,
		apiKey:     "testApiKey",
		httpClient: *mockServer.Client(),
	}

	data := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      "/test",
		Body:   nil,
	}

	var result struct {
		Result string `json:"result"`
	}

	err := client.ExecuteRequest(context.Background(), data, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Result != "success" {
		t.Errorf("Expected result 'success', got '%s'", result.Result)
	}
}

func TestExecuteRequest_SuccessfulPOST(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected 'POST', got '%s'", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected 'Content-Type' of 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"success"}`)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}))
	defer mockServer.Close()

	client := Client{
		httpClient: *mockServer.Client(),
		baseURL:    mustParseURL(mockServer.URL), // Helper function to parse URL
		apiKey:     "testApiKey",
	}

	payload := common.RequestPayload{
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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected 'POST', got '%s'", r.Method)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := io.WriteString(w, `{"message":"internal server error"}`)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}))
	defer mockServer.Close()

	client := Client{
		httpClient: *mockServer.Client(),
		baseURL:    mustParseURL(mockServer.URL),
		apiKey:     "testApiKey",
	}

	payload := common.RequestPayload{
		Method: "POST",
		U:      "/error",
		Body:   []byte(`{"key":"value"}`),
	}

	var response map[string]interface{} // Using interface{} to potentially capture any structure of response
	err := client.ExecuteRequest(context.Background(), payload, &response)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedErrorMessage := `processing response failed: {
	"message": "internal server error"
}`

	// Remove all whitespace when comparing the errors
	if strings.Join(strings.Fields(err.Error()), "") != strings.Join(strings.Fields(expectedErrorMessage), "") {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestAuthorizationKey(t *testing.T) {
	expectedAPIKey := "testApiKey"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+expectedAPIKey {
			t.Errorf("Authorization header does not match expected. Got %s, want %s", authHeader, "Bearer "+expectedAPIKey)
		}

		// If the Authorization header is as expected, respond with 200 OK
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"success"}`)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}))
	defer mockServer.Close()

	client := Client{
		httpClient: *mockServer.Client(),
		baseURL:    mustParseURL(mockServer.URL),
		apiKey:     expectedAPIKey,
	}

	payload := common.RequestPayload{
		Method: "GET",
		U:      "/",
	}

	var response interface{}

	err := client.ExecuteRequest(context.Background(), payload, &response)
	if err != nil {
		t.Fatalf("Did not expect an error, got %v", err)
	}
}

func mustParseURL(u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url
}

func TestIsScientificNotation(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expected    bool
	}{
		{
			description: "Valid scientific notation with positive exponent",
			input:       "1e+18",
			expected:    true,
		},
		{
			description: "Valid scientific notation with negative exponent",
			input:       "2.3e-4",
			expected:    true,
		},
		{
			description: "Invalid format - not scientific notation",
			input:       "not_scientific",
			expected:    false,
		},
		{
			description: "Invalid format - regular decimal number",
			input:       "3.14",
			expected:    false,
		},
		{
			description: "Valid scientific notation without sign in exponent",
			input:       "1e18",
			expected:    true,
		},
		{
			description: "Valid scientific notation with negative base and positive exponent",
			input:       "-3.5E+10",
			expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := isScientificNotation(tc.input)
			assert.Equal(t, tc.expected, result, fmt.Sprintf("TestIsScientificNotation failed for input: %s", tc.input))
		})
	}
}

func TestExpandScientificNotation(t *testing.T) {
	testCases := []struct {
		description    string
		input          string
		expectedOutput string
		expectedError  bool
	}{
		{
			description:    "Valid scientific notation with positive exponent",
			input:          "1e+18",
			expectedOutput: "1000000000000000000",
			expectedError:  false,
		},
		{
			description:    "Valid scientific notation with negative exponent",
			input:          "2.3e-4",
			expectedOutput: "0.00023",
			expectedError:  false,
		},
		{
			description:   "Invalid format - not a number",
			input:         "not_a_number",
			expectedError: true,
		},
		{
			description:    "Regular number without scientific notation",
			input:          "12345",
			expectedOutput: "12345",
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			output, err := expandScientificNotation(tc.input)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}
