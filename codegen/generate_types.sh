#!/bin/bash

# Global variables
verbose_logging=false
openapi_dir="openapi"
output_dir="generatedtypes"

# display_help shows the help message for this script
display_help() {
    echo "This script will generate request/response structs for all APIs supported by the SDK"
    echo
    echo "Usage: $0 [OPTIONS] [ARGUMENTS...]"
    echo ""
    echo "Options:"
    echo "  -logs              Enable verbose logging."
    echo "  -help              Display this help message."
}

remove_null_schemas() {
  local api_openapi_file_name="$1"
  local temp_file="${api_openapi_file_name}.tmp"
  jq 'walk(if type == "object" and has("anyOf")
          then .anyOf |= map(select(.type != "null"))
          else . end)' ${api_openapi_file_name} > ${temp_file}

  if [ $? -ne 0 ]; then
    echo "Error: Failed to remove null schemas with jq on $api_openapi_file_name."
    exit 1
  fi

  mv "${temp_file}" "${api_openapi_file_name}"
}

change_any_of_int_to_int() {
  local api_openapi_file_name="$1"
  local temp_file="${api_openapi_file_name}.tmp"

  jq '
    def modify_schema:
      if type == "object" and .schema.anyOf then
        .schema |= (.anyOf[0] + { description: .description, title: .title })
      else
        .
      end;

    .paths |= map_values(
      . as $path |
      . | map_values(
        if .parameters then
          .parameters |= map(modify_schema)
        else . end
      )
    )
  ' ${api_openapi_file_name} > ${temp_file}

  if [ $? -ne 0 ]; then
    echo "Error: Failed to change any of int schemas to int with jq on $api_openapi_file_name."
    exit 1
  fi

  mv "${temp_file}" "${api_openapi_file_name}"
}

change_any_of_ref_to_ref() {
  local api_openapi_file_name="$1"
  local temp_file="${api_openapi_file_name}.tmp"

  jq '
    def simplify_allOf:
      if type == "object" and .allOf then
        .allOf |= map(
          if type == "object" and has("$ref") then . else empty end
        ) |
        . |= if .allOf | length == 1 then .allOf[0] else . end |
        del(.allOf)
      else
        .
      end;

    # Apply the simplification to the paths
    .paths |= map_values(
      . as $path |
      . | map_values(
        if .parameters then
          .parameters |= map(
            if .schema then .schema |= simplify_allOf else . end
          )
        else . end |
        if .requestBody? then
          .requestBody.content."application/json".schema |= simplify_allOf
        else . end
      )
    ) |

    # Apply the simplification to the components schemas
    .components.schemas |= map_values(
      . |= simplify_allOf |
      if .properties then
        .properties |= map_values(simplify_allOf)
      else . end
    )
  ' ${api_openapi_file_name} > ${temp_file}

  if [ $? -ne 0 ]; then
    echo "Error: Failed to change any of ref schemas to ref with jq on $api_openapi_file_name."
    exit 1
  fi

  mv "${temp_file}" "${api_openapi_file_name}"
}

# update_operation_ids uses mappings from mapping.json to update operationId in the openapi files
update_operation_ids() {
    local api_openapi_file_name="$1"
    local mapping_file="mapping.json"  # Define the location of your mapping file
    local temp_file="${api_openapi_file_name}.tmp"  # Define the temporary file name

    # Ensure the mapping file exists
    if [ ! -f "$mapping_file" ]; then
        echo "Mapping file not found: $mapping_file"
        return 1
    fi

    # Use jq to update operationId based on the mapping file
    jq --slurpfile map "$mapping_file" '
        .paths |= with_entries(
            .value |= with_entries(
                if .value.operationId then
                    .value.operationId = ($map[0][.value.operationId] // .value.operationId)
                else
                    .
                end
            )
        )' "$api_openapi_file_name" > "$temp_file"

    if [ $? -ne 0 ]; then
      echo "Error: Failed to update operation ids with jq on $api_openapi_file_name."
      exit 1
    fi

    # Optionally, provide a message that the file has been updated
    echo "Updated operationIds in $api_openapi_file_name"
    mv "${temp_file}" "${api_openapi_file_name}"
}


# check_and_fix_incorrect_number_arrays checks for a known incorrect formatting for number arrays and fixes them for oapi-codegen
check_and_fix_incorrect_number_arrays() {
    local api_openapi_file_name="$1"

    if grep -q '"schema": { "type": "number\[\]" }' "$api_openapi_file_name"; then
        # If the pattern exists, use sed to replace the incorrect schema type
        echo "$(basename "$api_openapi_file_name") uses number arrays directly instead of using the array type. Fixing..."
        sed -i '' 's/"schema": { "type": "number\[\]" }/"schema": { "type": "array", "items": { "type": "number" } }/g' "$api_openapi_file_name" || {
            echo "Error while fixing number arrays in $api_openapi_file_name."
            exit 1
        }
    fi
}

# add_pointer_skip_field adds x-go-type-skip-optional-pointer to schema objects and parameters if not already present
# This is required to prevent the SDK from adding pointers to optional fields
add_pointer_skip_field() {
  local api_openapi_file_name="$1"
    local temp_file="${api_openapi_file_name}.tmp"  # Define the temporary file name

  jq '
           # Function to add x-go-type-skip-optional-pointer to schema objects if not already present
           def add_skip_pointer:
             if .type and (.["x-go-type-skip-optional-pointer"] // false) != true then
               . + {"x-go-type-skip-optional-pointer": true}
             else
               .
             end;

           # Apply to path parameters
           .paths |= map_values(
             . as $path |
             . | map_values(
               if .parameters then
                 .parameters |= map(
                   if .required == false and .schema then .schema |= add_skip_pointer else . end
                 )
               else . end |
               if .requestBody? then
                 .requestBody.content."application/json".schema |= add_skip_pointer
               else . end
             )
           ) |

           # Apply to components schemas
           .components.schemas |= map_values(
             if .properties then
               .properties |= map_values(add_skip_pointer)
             else . end
           )
         ' ${api_openapi_file_name} > ${temp_file}

  if [ $? -ne 0 ]; then
    echo "Error: Failed to add pointer skip fields with jq on $api_openapi_file_name."
    exit 1
  fi

  mv "${temp_file}" "${api_openapi_file_name}"
}

# Check that the script is being run from within the golang folder specifically
current_folder_name=$(basename "$PWD")

if [[ ! "$current_folder_name" == "codegen" ]]; then
    echo "This script can only be run from the directory in which it exists."
    exit 1
fi

# Check if Go is installed
if ! which go > /dev/null 2>&1; then
    echo "golang is not installed or not in PATH!"
    exit 1
fi

# Check if sed is installed
if ! which sed > /dev/null 2>&1; then
    echo "sed is not installed or not in PATH!"
    exit 1
fi

# Parse the command line arguments
index=1
while [ "$index" -le "$#" ]; do
    arg=${!index}
    case "$arg" in
        -logs)
            verbose_logging=true
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

# Install the type generator if it is not already installed
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.16.2 || {
    echo "Failed to install oapi-codegen."
    exit 1
}

# Check for any file in the openapi directory that doesn't fit the naming schema (all files must end with -openapi.json)
shopt -s nullglob # This ensures that the loop doesn't execute if no files match the pattern
for file in "$openapi_dir"/*; do
    filename=$(basename "$file")
    if [[ ! $filename =~ .+-openapi.json$ ]]; then
        echo "Error: file '$filename' does not match the expected naming schema of *-openapi.json."
        exit 1
    fi
done
shopt -u nullglob # Turn off the nullglob option

# Check if there are any openapi files to process
openapi_files_count=$(ls "$openapi_dir"/*-openapi.json 2>/dev/null | wc -l)
if [ "$openapi_files_count" -eq 0 ]; then
    echo "Warning: No openapi files found in $openapi_dir."
    exit 0
fi

# Loop over all openapi files in the directory
for api_openapi_file_name in "$openapi_dir"/*-openapi.json; do


    # Extract the service name from the filename
    package_name_raw=$(basename "$api_openapi_file_name" -openapi.json)

    types_file_name="${package_name_raw}_types.gen.go"

    # Check if package_name_raw contains an underscore
    if [[ $package_name_raw == *_* ]]; then
        # Remove underscore and everything after it
        package_name=${package_name_raw%%_*}
    else
        # If there is no underscore, assign raw value to package_name
        package_name=$package_name_raw
    fi

    if [ "$verbose_logging" == "true" ]; then
        echo "Generating type data for the $package_name directory"
        echo "Using the openapi file at $api_openapi_file_name"
        echo "The generated GoLang file will be in '$output_dir/$package_name/$types_file_name'"
        echo
    fi

    # Create a new directory with the service name under the output directory, if it doesn't exist
    mkdir -p "$output_dir/$package_name" || {
        echo "Error: Failed to create directory $output_dir/$package_name."
        exit 1
    }
    # Check for all known incorrect schema types and fix them if they exist
    check_and_fix_incorrect_number_arrays "$api_openapi_file_name"

    if [ "$verbose_logging" = "true" ]; then
      echo "Updating operationId for $api_openapi_file_name using mapping from $mapping_file"
    fi

    # Call to update operationIds
    update_operation_ids "$api_openapi_file_name"

    # Remove null schemas
    remove_null_schemas "$api_openapi_file_name"

    # simplify ints
    change_any_of_int_to_int "$api_openapi_file_name"

    # simplify ref
    change_any_of_ref_to_ref "$api_openapi_file_name"

    # Should always be the final schema step!
    # Add x-go-type-skip-optional-pointer to schema objects and parameters if not already present
    add_pointer_skip_field "$api_openapi_file_name"

    # Generate the oapi output into the new directory
    output_file="$output_dir/$package_name/${types_file_name}"
    oapi-codegen -generate types -package "$package_name" "$api_openapi_file_name" > "$output_file" || {
       echo "Error: Failed to generate types for $api_openapi_file_name."
       exit 1
    }

    # In the generated file, replace all 'form:' tags with 'url:' tags
    sed -i '' -E 's/`form:"([^"]+)"([^`]*)(`)/`url:"\1"\2\3/g' "$output_file" || {
        echo "Error: Failed to replace tags in $output_file."
        exit 1
    }

    new_dir="../sdk-clients/$package_name"  # Adjust relative path according to your project structure
        mkdir -p "$new_dir" && mv "$output_file" "$new_dir/" || {
            echo "Error: Failed to move generated types to $new_dir."
            exit 1
        }
done
