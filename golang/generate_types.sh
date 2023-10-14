#!/bin/bash

if [ "$#" -ge 1 ]; then
    # Check if the first argument is "logs"
    if [ "$1" == "logs" ]; then
        verbose_logging=true
    else
        echo "Error: Only 'logs' is accepted as an argument."
        exit 1
    fi
fi

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

fetch_and_generate() {
    local api_swagger_url="$1"
    local package_name="$2"
    local swagger_destination_directory="$3"

    local api_swagger_file_name="$package_name-swagger.json"

  if [ "$verbose_logging" ]; then
    echo "Generating type data for the $swagger_destination_directory directory"
    echo "Downloading the swagger file from $api_swagger_url"
    echo "The generated GoLang file will have a package name of '$package_name' and the file will be named '$api_swagger_file_name'"
  fi

    curl -s "$api_swagger_url" > "$api_swagger_file_name.tmp"

    # Check if curl was successful
    response_code=$?
    if [ $response_code -ne 0 ]; then
        echo "The curl command to get the latest swagger file from our servers (url: $api_swagger_url) failed with an error code of $response_code."
        exit 1
    fi

    # Check if the content contains "404"
    if grep -q '^{\"statusCode\":404' "$api_swagger_file_name.tmp"; then
        echo "The first contents of this file were a 404. The provided URL is likely wrong."
        echo "Manually check the contents of $api_swagger_file_name.tmp for more information."
        exit 2
    fi

    oapi-codegen -generate types -package "$package_name" "$api_swagger_file_name.tmp" > "$swagger_destination_directory/${package_name}_types.gen.go"

    rm "$api_swagger_file_name.tmp"
}

# Swap API
swap_api_swagger_url="https://api.1inch.io/swagger/ethereum-json"
swap_api_package_name="swap"
swap_api_swagger_destination_directory="client/swap"

fetch_and_generate "$swap_api_swagger_url" "$swap_api_package_name" "$swap_api_swagger_destination_directory"