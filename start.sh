#!/bin/bash

# Current directory
CUR_DIR=$(pwd)
SERVICE_DIR="$CUR_DIR/service"

# Function to build and execute a service
build_and_run_service() {
    local service_name="$1"
    local service_dir="$SERVICE_DIR/$service_name"

    if [ ! -d "$service_dir" ]; then
        echo "Service directory not found: $service_name"
        return
    fi

    if [ ! -f "$service_dir/build.sh" ]; then
        echo "build.sh not found for service: $service_name"
        return
    fi

    # Give execute permission to build.sh
    chmod +x "$service_dir/build.sh"
    # Run build.sh to generate output
    "$service_dir/build.sh"

    # Execute the generated executable
    if [ -f "$service_dir/output/bin/$service_name" ]; then
        "$service_dir/output/bin/$service_name"
    else
        echo "Executable not found for service: $service_name"
        return
    fi

    # Execute bootstrap.sh if it exists
    if [ -f "$service_dir/output/bootstrap.sh" ]; then
        "$service_dir/output/bootstrap.sh"
    fi
}

# Check if a specific service is provided as argument
if [ $# -gt 0 ]; then
    build_and_run_service "$1"
else
    # Run all services with build.sh
    for service in "$SERVICE_DIR"/*; do
        if [ -d "$service" ]; then
            service_name=$(basename "$service")
            if [ -f "$service/build.sh" ]; then
                build_and_run_service "$service_name"
            fi
        fi
    done
fi
