#!/bin/bash

# Test script for the Config API
# Make sure the server is running on :8080 before running this script

BASE_URL="http://localhost:8080"

echo "Testing Config API..."
echo "====================="

# Test 1: Create a config
echo "1. Creating a config..."
CREATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "database_config",
    "type": "database",
    "subtype": "postgresql",
    "tags": ["production", "database"],
    "metadata": {
      "host": "localhost",
      "port": 5432,
      "database": "mydb"
    }
  }')

# Extract HTTP status and response body
HTTP_STATUS=$(echo "$CREATE_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
RESPONSE_BODY=$(echo "$CREATE_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "HTTP Status: $HTTP_STATUS"
echo "Response: $RESPONSE_BODY"

if [ "$HTTP_STATUS" != "201" ]; then
    echo "Failed to create config. Status: $HTTP_STATUS"
    exit 1
fi

# Extract config ID from response using jq
CONFIG_ID=$(echo $RESPONSE_BODY | jq -r '.data.id // empty')
if [ -z "$CONFIG_ID" ] || [ "$CONFIG_ID" = "null" ]; then
    echo "Failed to create config or extract ID"
    echo "Response: $RESPONSE_BODY"
    exit 1
fi
echo "Config ID: $CONFIG_ID"

# Test 2: Get all configs
echo -e "\n2. Getting all configs..."
GET_ALL_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/configs?limit=10")
HTTP_STATUS=$(echo "$GET_ALL_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
RESPONSE_BODY=$(echo "$GET_ALL_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "HTTP Status: $HTTP_STATUS"
echo "$RESPONSE_BODY" | jq '.'

# Test 3: Get config by ID
if [ ! -z "$CONFIG_ID" ]; then
    echo -e "\n3. Getting config by ID: $CONFIG_ID"
    GET_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/config/$CONFIG_ID")
    HTTP_STATUS=$(echo "$GET_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
    RESPONSE_BODY=$(echo "$GET_RESPONSE" | sed '/HTTP_STATUS:/d')
    
    echo "HTTP Status: $HTTP_STATUS"
    echo "$RESPONSE_BODY" | jq '.'
fi

# Test 4: Update config
if [ ! -z "$CONFIG_ID" ]; then
    echo -e "\n4. Updating config: $CONFIG_ID"
    UPDATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/config/$CONFIG_ID" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "updated_database_config",
        "subtype": "mysql",
        "tags": ["staging", "database"],
        "metadata": {
          "host": "staging.example.com",
          "port": 3306,
          "database": "staging_db"
        }
      }')
    HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
    RESPONSE_BODY=$(echo "$UPDATE_RESPONSE" | sed '/HTTP_STATUS:/d')
    
    echo "HTTP Status: $HTTP_STATUS"
    echo "$RESPONSE_BODY" | jq '.'
fi

# Test 5: Delete config
if [ ! -z "$CONFIG_ID" ]; then
    echo -e "\n5. Deleting config: $CONFIG_ID"
    DELETE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X DELETE "$BASE_URL/config/$CONFIG_ID")
    HTTP_STATUS=$(echo "$DELETE_RESPONSE" | grep "HTTP_STATUS:" | cut -d':' -f2)
    RESPONSE_BODY=$(echo "$DELETE_RESPONSE" | sed '/HTTP_STATUS:/d')
    
    echo "HTTP Status: $HTTP_STATUS"
    if [ "$HTTP_STATUS" = "200" ]; then
        echo "Config deleted successfully!"
    else
        echo "Failed to delete config. Response: $RESPONSE_BODY"
    fi
fi

echo -e "\nAPI testing completed!"
