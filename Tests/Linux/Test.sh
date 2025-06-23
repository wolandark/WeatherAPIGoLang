#!/usr/bin/env bash

BASE_URL="http://localhost:8080"
WEATHER_ID=""

# Colors, Because I can
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

print_color() {
    local color=$1
    local message=$2
    printf "${color}%s${NC}\n" "$message"
}

# Check If jq Is Installed
check_jq() {
    if ! command -v jq >/dev/null 2>&1; then
        print_color "$RED" "Error: jq is required but not installed. Please install jq first."
        exit 1
    fi
}

# Check If Curl Is Installed
check_curl() {
    if ! command -v curl >/dev/null 2>&1; then
        print_color "$RED" "Error: curl is required but not installed. Please install curl first."
        exit 1
    fi
}

# Check Deps
check_jq
check_curl

print_color "$GREEN" "=== Weather API Test Script ==="
echo ""

# Create weather
print_color "$YELLOW" "1. Testing POST /weather - Creating weather record for London, UK"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/weather" \
    -H "Content-Type: application/json" \
    -d '{"cityName": "London", "country": "UK"}')

print_color "$CYAN" "Response: $CREATE_RESPONSE"

# Extract ID from JSON response using jq
if [ -n "$CREATE_RESPONSE" ]; then
    WEATHER_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id // empty')
    if [ -n "$WEATHER_ID" ] && [ "$WEATHER_ID" != "null" ]; then
        print_color "$GREEN" "Extracted Weather ID: $WEATHER_ID"
    else
        print_color "$RED" "Failed to extract ID from response"
        exit 1
    fi
else
    print_color "$RED" "Failed to create weather record"
    exit 1
fi

echo ""

# Get all
print_color "$YELLOW" "2. Testing GET /weather - Fetching all weather records"
GET_ALL_RESPONSE=$(curl -s "$BASE_URL/weather")
print_color "$CYAN" "Response: $GET_ALL_RESPONSE"
echo ""

# Get weather by ID
print_color "$YELLOW" "3. Testing GET /weather/:id - Fetching weather record by ID: $WEATHER_ID"
GET_BY_ID_RESPONSE=$(curl -s "$BASE_URL/weather/$WEATHER_ID")
print_color "$CYAN" "Response: $GET_BY_ID_RESPONSE"
echo ""

# Update
print_color "$YELLOW" "4. Testing PUT /weather/:id - Updating weather record with mock data"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/weather/$WEATHER_ID" \
    -H "Content-Type: application/json" \
    -d '{"temperature": 25.5, "description": "sunny", "humidity": 65}')
print_color "$CYAN" "Response: $UPDATE_RESPONSE"
echo ""

# Get latest weather by city
print_color "$YELLOW" "5. Testing GET /weather/latest/:cityName - Fetching latest weather for London"
GET_LATEST_RESPONSE=$(curl -s "$BASE_URL/weather/latest/London")
print_color "$CYAN" "Response: $GET_LATEST_RESPONSE"
echo ""

# Delete
print_color "$YELLOW" "6. Testing DELETE /weather/:id - Deleting weather record ID: $WEATHER_ID"
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/weather/$WEATHER_ID")
print_color "$CYAN" "Response: $DELETE_RESPONSE"
echo ""

# Verify deletion - try get deleted record
print_color "$YELLOW" "7. Verifying deletion - Attempting to fetch deleted record"
VERIFY_DELETE_RESPONSE=$(curl -s "$BASE_URL/weather/$WEATHER_ID")
print_color "$CYAN" "Response: $VERIFY_DELETE_RESPONSE"
echo ""

print_color "$GREEN" "=== Test Script Completed ==="
print_color "$GREEN" "All CRUD operations have been tested."

# Optionally Create a few more test records with different cities
echo ""
print_color "$MAGENTA" "=== Creating Additional Test Records ==="

CITIES='[
    {"city": "Paris", "country": "France"},
    {"city": "Tokyo", "country": "Japan"},
    {"city": "New York", "country": "US"}
]'

echo "$CITIES" | jq -c '.[]' | while read -r location; do
    CITY=$(echo "$location" | jq -r '.city')
    COUNTRY=$(echo "$location" | jq -r '.country')
    
    print_color "$YELLOW" "Creating weather record for $CITY, $COUNTRY"
    
    TEST_RESPONSE=$(curl -s -X POST "$BASE_URL/weather" \
        -H "Content-Type: application/json" \
        -d "{\"cityName\": \"$CITY\", \"country\": \"$COUNTRY\"}")
    
    print_color "$CYAN" "Response: $TEST_RESPONSE"
    # Dont Overwhelm the api 💅
    sleep 1
done

echo ""
print_color "$YELLOW" "Final GET all records:"
FINAL_GET_ALL=$(curl -s "$BASE_URL/weather")
print_color "$CYAN" "Response: $FINAL_GET_ALL"

echo ""
print_color "$GREEN" "=== All Tests Completed Successfully ==="

# Optional: Pretty print final results with jq
echo ""
print_color "$MAGENTA" "=== Summary of Final Records ==="
if [ -n "$FINAL_GET_ALL" ]; then
    echo "$FINAL_GET_ALL" | jq -r '.[] | "ID: \(.id) | City: \(.cityName), \(.country) | Temp: \(.temperature)°C | \(.description)"'
else
    print_color "$RED" "No records found or error retrieving records"
fi

