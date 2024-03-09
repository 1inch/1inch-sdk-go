This SDK will be open for contributions from the community. Contribution guidelines will be added soon!


### Versioning

This library is currently in the developer preview phase (versions 0.x.x). There will be significant changes to the design of this library leading up to a 1.0.0 release. You can expect the API calls, library structure, etc. to break between each release. Once the library version reaches 1.0.0 and beyond, it will follow traditional semver conventions.

### Project structure

This SDK is powered by a [client struct](https://github.com/1inch/1inch-sdk/blob/main/golang/client/client.go) that contains instances of all Services used to talk to the 1inch APIs

Each Service maps 1-to-1 with the underlying Dev Portal REST API. See [SwapService](https://github.com/1inch/1inch-sdk/blob/main/golang/client/swap.go) as an example. Under each function, you will find the matching REST API path)

Each Service uses various types and functions to do its job that are kept separate from the main service file. These can be found in the accompanying folder within the client directory (see the [swap](https://github.com/1inch/1inch-sdk/tree/main/golang/client/swap) package)

### Type generation

Type generation is done using the `generate_types.sh` script. To add a new swagger file or update an existing one, place the swagger file in `swagger-static` and run the script. It will generate the types file and place it in the appropriately-named sub-folder inside the `client` directory

### Swagger file formatting
For consistency, Swagger files should be formatted with `prettier`

This can be installed globally using npm:

`npm install -g prettier`

If using GoLand, you can set up this action to run automatically using File Watchers:

1. Go to Settings or Preferences > Tools > File Watchers.
2. Click the + button to add a new watcher.
3. For `File type`, choose JSON.
4. For `Scope`, choose Project Files.
5. For `Program`, provide the path to the `prettier`. This can be gotten by running `which prettier`.
6. For `Arguments`, use `--write $FilePath$`.
7. For `Output paths to refresh`, use `$FilePath$`.
8. Ensure the Auto-save edited files to trigger the watcher option is checked
