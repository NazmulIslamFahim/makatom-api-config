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

echo "Environment variables set:"
echo "  API_PORT: $API_PORT"
echo "  ENVIRONMENT: $ENVIRONMENT"
echo "  DEBUG: $DEBUG"
echo "  MONGO_URI: $MONGO_URI"
echo "  MONGO_DATABASE: $MONGO_DATABASE"

echo ""
echo "Starting server..."
echo "Press Ctrl+C to stop"

# Run the server
go run cmd/main.go
