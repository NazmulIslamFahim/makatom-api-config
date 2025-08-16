#!/bin/bash

# Test script for the Config API
# Make sure the server is running on :8080 before running this script

BASE_URL="http://localhost:8080"

echo "Testing Config API..."
echo "====================="

# Test 1: Create a config
echo "1. Creating a config..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "database_config",
    "type": "database",
    "subtype": "postgresql",
    "tags": ["production", "database"],
    "tenant_id": "tenant123",
    "metadata": {
      "host": "localhost",
      "port": 5432,
      "database": "mydb"
    }
  }')

echo "Create Response: $CREATE_RESPONSE"
CONFIG_ID=$(echo $CREATE_RESPONSE | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
echo "Config ID: $CONFIG_ID"

# Test 2: Get all configs
echo -e "\n2. Getting all configs..."
curl -s -X GET "$BASE_URL/configs?tenant_id=tenant123&limit=10" | jq '.'

# Test 3: Get config by ID
if [ ! -z "$CONFIG_ID" ]; then
    echo -e "\n3. Getting config by ID: $CONFIG_ID"
    curl -s -X GET "$BASE_URL/config/get?id=$CONFIG_ID" | jq '.'
fi

# Test 4: Update config
if [ ! -z "$CONFIG_ID" ]; then
    echo -e "\n4. Updating config: $CONFIG_ID"
    curl -s -X PUT "$BASE_URL/config/update?id=$CONFIG_ID" \
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
      }' | jq '.'
fi

# Test 5: Delete config
if [ ! -z "$CONFIG_ID" ]; then
    echo -e "\n5. Deleting config: $CONFIG_ID"
    curl -s -X DELETE "$BASE_URL/config/delete?id=$CONFIG_ID"
    echo -e "\nConfig deleted successfully!"
fi

echo -e "\nAPI testing completed!"
