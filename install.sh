#!/bin/bash

# Define the binary name
BINARY_NAME="wakey"

# Build the Go binary
GOOS=darwin GOARCH=amd64 go build -o $BINARY_NAME

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

# Move the binary to /usr/local/bin
mv $BINARY_NAME /usr/local/bin/

# Verify if the binary was moved successfully
if [ $? -eq 0 ]; then
    echo "$BINARY_NAME has been successfully installed in /usr/local/bin."
else
    echo "Failed to move $BINARY_NAME to /usr/local/bin."
    exit 1
fi
