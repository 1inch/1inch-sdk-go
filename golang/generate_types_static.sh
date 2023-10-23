#!/bin/bash

display_help() {
    echo "This script will generate request/response structs for all APIs supported by the SDK"
    echo
    echo "Usage: $0 [OPTIONS] [ARGUMENTS...]"
    echo ""
    echo "Options:"
    echo "  -logs              Enable verbose logging."
    echo "  -swagger-dir <path> Specify the directory containing swagger files. For each file named '<service>-swagger.json',"
    echo "                      code will be generated inside a directory named '<service>'."
    echo "  -output-dir <path>  Specify the directory where generated code should be placed."
    echo "  -help              Display this help message."
}

check_and_fix_incorrect_number_arrays() {
    local api_swagger_file_name="$1"

    if grep -q '"schema": { "type": "number\[\]" }' "$api_swagger_file_name"; then
        # If the pattern exists, use sed to replace the incorrect schema type
        echo "$(basename "$api_swagger_file_name") uses number arrays directly instead of using the array type. Fixing..."
        sed -i.bak 's/"schema": { "type": "number\[\]" }/"schema": { "type": "array", "items": { "type": "number" } }/g' "$api_swagger_file_name" || {
            echo "Error while fixing number arrays in $api_swagger_file_name."
            exit 1
        }
        rm "$api_swagger_file_name.bak"
    fi
}

# Check if Go is installed
if ! which go > /dev/null 2>&1; then
    echo "Golang is not installed or not in PATH!"
    exit 1
fi

# Check for sed
if ! which sed > /dev/null 2>&1; then
    echo "sed is not installed or not in PATH!"
    exit 1
fi

verbose_logging=false
swagger_dir=""
output_dir=""

index=1
while [ "$index" -le "$#" ]; do
    arg=${!index}
    case "$arg" in
        -logs)
            verbose_logging=true
            ;;
        -swagger-dir)
            index=$((index + 1))
            if [ "$index" -le "$#" ]; then
                swagger_dir=${!index}
            else
                echo "Error: -swagger-dir flag requires a directory value."
                exit 1
            fi
            ;;
        -output-dir)
            index=$((index + 1))
            if [ "$index" -le "$#" ]; then
                output_dir=${!index}
            else
                echo "Error: -output-dir flag requires a directory value."
                exit 1
            fi
            ;;
        -help)
            display_help
            ;;
        *)
            echo "Error: Unknown argument '$arg'."
            echo "Use '-help' for a list of available options."
            echo
            display_help
            exit 1
            ;;
    esac
    index=$((index + 1))
done

# Check if either -swagger-dir or -output-dir is missing
if [ -z "$swagger_dir" ]; then
    echo "Error: Missing argument -swagger-dir. Please specify the directory containing swagger files."
    exit 1
fi

if [ -z "$output_dir" ]; then
    echo "Error: Missing argument -output-dir. Please specify the directory where the generated code should be placed."
    exit 1
fi

# Install the type generator if it is not already installed
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest || {
    echo "Failed to install oapi-codegen."
    exit 1
}

# Check if the swagger directory exists and is readable
if [ ! -d "$swagger_dir" ] || [ ! -r "$swagger_dir" ]; then
    echo "Error: Swagger directory $swagger_dir not found or not readable."
    exit 1
fi

# Check for any file in the directory that doesn't fit the naming schema
for file in "$swagger_dir"/*; do
    filename=$(basename "$file")
    if [[ ! $filename =~ .+-swagger.json$ ]]; then
        echo "Error: file '$filename' does not match the expected naming schema of *-swagger.json."
        exit 1
    fi
done

# Check if there are any swagger files to process
swagger_files_count=$(ls "$swagger_dir"/*-swagger.json 2>/dev/null | wc -l)
if [ "$swagger_files_count" -eq 0 ]; then
    echo "Warning: No swagger files found in $swagger_dir."
    exit 0
fi

# Loop over all swagger files in the directory
for api_swagger_file_name in "$swagger_dir"/*-swagger.json; do
    # Extract the service name from the filename
    package_name=$(basename "$api_swagger_file_name" -swagger.json)

    types_file_name="${package_name}_types.gen.go"

    if [ "$verbose_logging" == "true" ]; then
        echo "Generating type data for the $package_name directory"
        echo "Using the swagger file at $api_swagger_file_name"
        echo "The generated GoLang file will be in '$output_dir/$package_name/$types_file_name'"
        echo
    fi

    # Create a new directory with the service name under the output directory, if it doesn't exist
    mkdir -p "$output_dir/$package_name" || {
        echo "Error: Failed to create directory $output_dir/$package_name."
        exit 1
    }

    # Check for all known incorrect schema types and fix them if they exist
    check_and_fix_incorrect_number_arrays "$api_swagger_file_name"

    # Generate the swagger output into the new directory
    output_file="$output_dir/$package_name/${types_file_name}"
    oapi-codegen -generate types -package "$package_name" "$api_swagger_file_name" > "$output_file" || {
       echo "Error: Failed to generate types for $api_swagger_file_name."
       exit 1
    }

    # In the generated file, replace all 'form:' tags with 'url:' tags
    sed -i.bak -E 's/`form:"([^"]+)"([^`]*)(`)/`url:"\1"\2\3/g' "$output_file" || {
        echo "Error: Failed to replace tags in $output_file."
        exit 1
    }
    rm "$output_file.bak"
done
