This SDK will be open for contributions from the community. Contribution guidelines will be added soon!

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
