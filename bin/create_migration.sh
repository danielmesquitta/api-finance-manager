#!/bin/bash

directory="$1"
migration_name="$2"

# Check if a migration name is provided
if [ "$#" -ne 2 ]; then
  echo "Error: Migration name is missing."
  if [ "$directory" == "sql/migrations" ]; then
    echo "Usage: make create_migration <migration_name>"
    exit 1
  fi

  if [ "$directory" == "sql/testdata" ]; then
    echo "Usage: make create_testdata <migration_name>"
    exit 1
  fi

  echo "Usage: $0 <directory> <migration_name>"
  exit 1
fi

# Get the current UTC timestamp formatted as YYYYMMDDHHMMSS
timestamp=$(date -u +'%Y%m%d%H%M%S')

# Define the target filename
filename="${timestamp}_${migration_name}.sql"

# Create the directory if it doesn't exist
mkdir -p "$directory"

# Create the file in the target directory
touch "${directory}/${filename}"

echo "Created migration file: ${directory}/${filename}"
