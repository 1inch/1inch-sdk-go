package balances

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/1inch/1inch-sdk-go/common"
)

func TestNewClient(t *testing.T) {
	//mockAPI := api{
	//	chainId: 1,
	//}

}

type mockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj interface{}
}

func (m *mockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v interface{}) error {
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

//var mockedSwapHttpApiResp = {}
