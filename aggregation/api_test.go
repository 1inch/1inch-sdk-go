package aggregation

import (
	"context"
	"fmt"
	"github.com/1inch/1inch-sdk-go/internal/common"
	"reflect"
	"testing"
)

// MockHttpExecutor is a mock implementation of the httpExecutor interface
type MockHttpExecutor struct {
	// Add fields to control behavior or track method calls during tests
	Called      bool
	ExecuteErr  error
	ResponseObj interface{}
}

// ExecuteRequest simulates the behavior of an HTTP request execution
func (m *MockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v interface{}) error {
	m.Called = true
	if m.ExecuteErr != nil {
		return m.ExecuteErr
	}

	// Copy the mock response object to v
	if m.ResponseObj != nil && v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("v must be a non-nil pointer")
		}
		reflect.Indirect(rv).Set(reflect.ValueOf(m.ResponseObj))
	}
	return nil
}

// TestGetQuote tests the GetQuote function of the api struct
func TestGetQuote(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockExecutor := &MockHttpExecutor{
		ResponseObj: QuoteResponse{ /* initialize this with expected data */ },
	}
	api := api{httpExecutor: mockExecutor}

	// Define test params
	params := GetQuoteParams{ /* initialize with valid parameters */ }

	// Execute the test function
	quote, err := api.GetQuote(ctx, params)
	if err != nil {
		t.Fatalf("GetQuote returned an error: %v", err)
	}

	// Validate the mockExecutor was called
	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	// Further validations can be added here to check if 'quote' contains the expected data

	// Example: Check if the quote contains expected data
	expectedQuote := QuoteResponse{ /* define expected quote response */ }
	if !reflect.DeepEqual(*quote, expectedQuote) {
		t.Errorf("Expected quote to be %+v, got %+v", expectedQuote, *quote)
	}
}

// temp here
// todo: remove and rewrite
//func TestGetQuoteReal(t *testing.T) {
//	ctx := context.Background()
//
//	base, err := url.Parse("https://api.1inch.dev/")
//	if err != nil {
//		t.Fatalf("failed with an URL base: %v", err)
//	}
//	client := http_executor.DefaultHttpClient(base, "YOUR API KEY")
//
//	a := api{
//		httpExecutor: &client,
//	}
//
//	resp, err := a.GetQuote(ctx, GetQuoteParams{
//		ChainId: chains.Ethereum,
//		AggregationControllerGetQuoteParams: AggregationControllerGetQuoteParams{
//			Src:    "0x6b175474e89094c44da98b954eedeac495271d0f",
//			Dst:    "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//			Amount: "1000000000000000000",
//		},
//	})
//
//	if err != nil {
//		t.Fatalf("GetQuote returned an error: %v", err)
//	}
//	expectedQuote := QuoteResponse{ /* define expected quote response */ }
//	if !reflect.DeepEqual(*resp, expectedQuote) {
//		t.Errorf("Expected quote to be %+v, got %+v", expectedQuote, *resp)
//	}
//}
