#!/bin/bash

# Check if Go is installed
if ! which go > /dev/null 2>&1; then
    echo "Golang is not installed!"
    exit 1
fi

verbose_logging=false
skip_downloads=false

# Parse all CLI flags
for arg in "$@"; do
    case "$arg" in
        -logs)
            verbose_logging=true
            ;;
        -skip-downloads)
            skip_downloads=true
            ;;
        -help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -logs            Enable verbose logging."
            echo "  -skip-downloads  Skip the download processes."
            echo "  -help            Display this help message."
            exit 0
            ;;
        *)
            echo "Error: Unknown argument '$arg'."
            echo "Use '-help' for a list of available options."
            exit 1
            ;;
    esac
done

# Install the type generator if it is not already installed
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

fetch_and_generate() {
    local api_swagger_url="$1"
    local package_name="$2"
    local types_destination_directory="$3"

    local api_swagger_file_name="swagger/$package_name-swagger.json"

  if [ "$verbose_logging" == "true" ]; then
    echo "Generating type data for the $types_destination_directory directory"
    echo "Downloading the swagger file from $api_swagger_url"
    echo "The generated GoLang file will have a package name of '$package_name' and the file will be at '$api_swagger_file_name'"
    echo
  fi

  if [ "$skip_downloads" == "false" ]; then
      curl -s "$api_swagger_url" > "$api_swagger_file_name"

      # Check if curl was successful
      response_code=$?
      if [ $response_code -ne 0 ]; then
          echo "The curl command to get the latest swagger file from our servers (url: $api_swagger_url) failed with an error code of $response_code."
          exit 1
      fi

      # Check if the content contains "404"
      if grep -q '^{\"statusCode\":404' "$api_swagger_file_name"; then
          echo "The first contents of $api_swagger_file_name were a 404. The provided URL is likely wrong."
          echo "Manually check the contents of $api_swagger_file_name for more information."
          exit 2
      fi

      # Check if the content contains data that looks like an openapi spec
      if ! grep -q '^{\"openapi\"' "$api_swagger_file_name"; then
          echo "The first contents of $api_swagger_file_name does not look like openapi data. The request has likely failed."
          echo "Manually check the contents of $api_swagger_file_name for more information."
          exit 3
      fi
  fi

    oapi-codegen -generate types -package "$package_name" "$api_swagger_file_name" > "$types_destination_directory/${package_name}_types.gen.go"
}

# Swap API
# TODO This URL is not versioned. We will want versioned APIs in the future.
fetch_and_generate "https://api.1inch.io/swagger/ethereum-json" "swap" "client/swap"

# Spot Price API
fetch_and_generate "https://token-prices.1inch.io/swagger/ethereum-json" "spotprice" "client/spotprice"