#!/bin/sh

# Source file path
source_file="$(ls /vault-plugin-*)"

# Destination directory path
destination="/plugins/$(ls vault-plugin-*)"

# Maximum wait time in seconds
max_wait_time=60

# Start time
start_time=$(date +%s)

# Copy the file
# Wait for directory /plugins to exist or until max wait time is reached
while [ ! -d "/plugins" ] && [ $(( $(date +%s) - start_time )) -lt $max_wait_time ]; do
    sleep 1
done

# Check if /plugins directory exists
if [ -d "/plugins" ]; then
    cp "$source_file" "$destination"
else
    echo "Max wait time of $max_wait_time seconds reached. /plugins directory not found."
fi