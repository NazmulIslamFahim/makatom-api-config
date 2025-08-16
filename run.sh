#!/bin/bash

# Run script for the Config API
# This script sets up the environment and runs the server

echo "Starting Makatom API Config Service..."

# Set environment variables
export API_PORT=":8080"
export ENVIRONMENT="development"
export DEBUG="true"
export MONGO_URI="mongodb://localhost:27017/makatom_config"
export MONGO_DATABASE="makatom_config"
export MONGO_URI_NAME="config"

echo "Environment variables set:"
echo "  API_PORT: $API_PORT"
echo "  ENVIRONMENT: $ENVIRONMENT"
echo "  DEBUG: $DEBUG"
echo "  MONGO_URI: $MONGO_URI"
echo "  MONGO_DATABASE: $MONGO_DATABASE"
echo "  MONGO_URI_NAME: $MONGO_URI_NAME"

echo ""
echo "Starting server..."

# Check if air is available for hot reloading
if command -v air &> /dev/null; then
    echo "Using Air for hot reloading..."
    echo "Press Ctrl+C to stop"
    air
else
    echo "Air not found, using regular go run..."
    echo "Press Ctrl+C to stop"
    go run cmd/main.go
fi
