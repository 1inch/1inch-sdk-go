package client

//
//import (
//	"fmt"
//	"net/http"
//	"net/http/httptest"
//	"net/url"
//	"os"
//
//	"github.com/1inch/1inch-sdk-go/internal/helpers/constants/chains"
//
//	"github.com/1inch/1inch-sdk-go/client/models"
//)
//
//// setup sets up a test HTTP server along with a Client that is
//// configured to talk to that test server. Tests should register handlers on
//// mux which provide mock responses for the API method being tested.
//func setup() (*Client, *http.ServeMux, string, func(), error) {
//	mux := http.NewServeMux()
//
//	// This defaults all requests to return a 404
//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusNotFound)
//		fmt.Fprintf(w, "Not Found")
//	})
//
//	// server is a test HTTP server used to provide mock API responses.
//	// the base URL of the client will have its destination swapped to use this new test server for requests
//	server := httptest.NewServer(mux)
//	c, err := NewClient(
//		models.ClientConfig{
//			DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
//			Web3HttpProviders: []models.Web3Provider{
//				{
//					ChainId: chains.Ethereum,
//					Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
//				},
//			},
//		})
//	if err != nil {
//		return nil, nil, "", nil, err
//	}
//
//	url, _ := url.Parse(server.URL + "/")
//	c.ApiBaseURL = url
//	return c, mux, server.URL, server.Close, nil
//}
