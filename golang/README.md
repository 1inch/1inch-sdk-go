# Dev Portal Go SDK

## Running locally

To authenticate your client with the 1inch API, your API token will be needed by the initial config used to create the SDK client. It is recommended to store this information in your local environment and read it dynamically at runtime:

```go
...
config := client.Config{
    ApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
}

c, err := client.NewClient(config)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
...
```