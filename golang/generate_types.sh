#!/bin/bash

# Constants for exit codes
GOLANG_NOT_INSTALLED=1
CURL_NOT_INSTALLED=2
UNKNOWN_ARGUMENT=3
CURL_FAIL=4
SWAGGER_404=5
SWAGGER_INVALID_DATA=6

# Check if Go is installed
if ! which go > /dev/null 2>&1; then
    echo "Golang is not installed!"
    exit $GOLANG_NOT_INSTALLED
fi

# Check if curl is installed
if ! which curl > /dev/null 2>&1; then
    echo "curl is not installed!"
    exit $CURL_NOT_INSTALLED
fi

display_help() {
    echo "This script will generate request/response structs for all APIs supported by the SDK"
    echo
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -logs            Enable verbose logging."
    echo "  -use-local       Use the swagger files in swagger-static."
    echo "  -help            Display this help message."
    exit 0
}

verbose_logging=false
use_local=false

# Parse all CLI flags
for arg in "$@"; do
    case "$arg" in
        -logs)
            verbose_logging=true
            ;;
        -use-local)
            use_local=true
            ;;
        -help)
            display_help
            ;;
        *)
            echo "Error: Unknown argument '$arg'."
            echo "Use '-help' for a list of available options."
            echo
            display_help
            exit $UNKNOWN_ARGUMENT
            ;;
    esac
done

# Install the type generator if it is not already installed
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

fetch_and_generate() {
    local api_swagger_url="$1"
    local package_name="$2"
    local types_destination_directory="$3"

    local types_file_name="${package_name}_types.gen.go"

    if [ "$verbose_logging" == "true" ]; then
        echo "Generating type data for the $types_destination_directory directory"
        echo "Downloading the swagger file from $api_swagger_url"
        echo "The generated GoLang file will be at '$types_destination_directory/$types_file_name'"
        echo
    fi

    local api_swagger_file_name="swagger-static/$package_name-swagger.json"

    if [ "$use_local" == "false" ]; then

        api_swagger_file_name="swagger-dynamic/$package_name-swagger.json"
        curl -s "$api_swagger_url" > "$api_swagger_file_name"

        # Check if curl was successful
        response_code=$?
        if [ $response_code -ne 0 ]; then
            echo "The curl command to get the latest swagger file from our servers (url: $api_swagger_url) failed with an error code of $response_code."
            exit $CURL_FAIL
        fi

        # Check if the content contains "404"
        if grep -q '^{\"statusCode\":404' "$api_swagger_file_name"; then
            echo "The first contents of $api_swagger_file_name were a 404. The provided URL is likely wrong."
            echo "Manually check the contents of $api_swagger_file_name for more information."
            exit $SWAGGER_404
        fi

        # Check if the content contains data that looks like an openapi spec
        if ! grep -q '^{\"openapi\"' "$api_swagger_file_name"; then
            echo "The first contents of $api_swagger_file_name does not look like openapi data. The request has likely failed."
            echo "Manually check the contents of $api_swagger_file_name for more information."
            exit $SWAGGER_INVALID_DATA
        fi
    fi

    local output_file=$types_destination_directory/${types_file_name}
    oapi-codegen -generate types -package "$package_name" "$api_swagger_file_name" > "$output_file"

    # In the generated file, replace all 'form:' tags with 'url:' tags so the fields are properly parsed as URL query parameters
    sed -i.bak -E 's/`form:"([^"]+)"([^`]*)(`)/`url:"\1"\2\3/g' "$output_file"

    # Remove the backup file created by sed
    rm "$output_file.bak"
}



# Swap API
# TODO This URL is not versioned. We will want versioned APIs in the future.
fetch_and_generate "https://api.1inch.io/swagger/ethereum-json" "swap" "client/swap"

# Spot Price API
fetch_and_generate "https://token-prices.1inch.io/swagger/ethereum-json" "tokenprices" "client/tokenprices"