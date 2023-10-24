# Dev Portal Go SDK

## Using the SDK in your project

### Client Initialization
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

### API calls

Once you have an initialized client, API calls can be made using any of the client's underlying services:

// TODO add simple example here and link to `examples` directory 

## Development

### Type generation

// TODO explain this more

Type generation is done using the `generate_types.sh` script 

### Swagger file formatting
For consistency, Swagger files should be formatted with `prettier`

This can be installed globally using npm:

`npm install -g prettier`

If using GoLand, you can setup this action to run automatically using File Watchers:

1. Go to Settings or Preferences > Tools > File Watchers.
2. Click the + button to add a new watcher.
3. For `File type`, choose JSON.
4. For `Scope`, choose Project Files.
5. For `Program`, provide the path to the `prettier`. This can be gotten by running `which prettier`.
6. For `Arguments`, use `--write $FilePath$`.
7. For `Output paths to refresh`, use `$FilePath$`.
8. Ensure the Auto-save edited files to trigger the watcher option is checked
