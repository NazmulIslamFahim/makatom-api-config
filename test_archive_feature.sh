#!/bin/bash

# Test script for archive feature
BASE_URL="http://localhost:8080"

echo "=== Testing Archive Feature ==="
echo

# Test 1: Create a config
echo "1. Creating a config..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-config-archive",
    "type": "database",
    "subtype": "postgres",
    "tags": ["production", "critical"],
    "metadata": {
      "host": "localhost",
      "port": 5432,
      "database": "testdb"
    }
  }')

echo "Create Response: $CREATE_RESPONSE"
CONFIG_ID=$(echo $CREATE_RESPONSE | grep -o '"_id":"[^"]*"' | cut -d'"' -f4)
echo "Config ID: $CONFIG_ID"
echo

# Test 2: Update the config (this should create an archive)
echo "2. Updating the config (should create archive)..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/config?id=$CONFIG_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-config-archive-updated",
    "metadata": {
      "host": "localhost",
      "port": 5432,
      "database": "testdb",
      "ssl": true
    }
  }')

echo "Update Response: $UPDATE_RESPONSE"
echo

# Test 3: Update again (should create another archive)
echo "3. Updating the config again (should create another archive)..."
UPDATE_RESPONSE2=$(curl -s -X PUT "$BASE_URL/config?id=$CONFIG_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "tags": ["production", "critical", "updated"],
    "metadata": {
      "host": "localhost",
      "port": 5432,
      "database": "testdb",
      "ssl": true,
      "max_connections": 100
    }
  }')

echo "Second Update Response: $UPDATE_RESPONSE2"
echo

# Test 4: Get config archives
echo "4. Getting config archives..."
ARCHIVES_RESPONSE=$(curl -s -X GET "$BASE_URL/config/archives?id=$CONFIG_ID")
echo "Archives Response: $ARCHIVES_RESPONSE"
echo

# Test 5: Get the current config
echo "5. Getting current config..."
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/config?id=$CONFIG_ID")
echo "Get Response: $GET_RESPONSE"
echo

# Test 6: Delete the config (should also delete all archives)
echo "6. Deleting the config (should also delete all archives)..."
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/config?id=$CONFIG_ID")
echo "Delete Response: $DELETE_RESPONSE"
echo

# Test 7: Try to get archives after deletion (should return 404)
echo "7. Trying to get archives after deletion (should return 404)..."
ARCHIVES_AFTER_DELETE_RESPONSE=$(curl -s -X GET "$BASE_URL/config/archives?id=$CONFIG_ID")
echo "Archives After Delete Response: $ARCHIVES_AFTER_DELETE_RESPONSE"
echo

echo "=== Archive Feature Test Complete ==="
