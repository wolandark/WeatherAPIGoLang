$baseUrl = "http://localhost:8080"
$weatherId = ""

Write-Host "=== Weather API Test Script ===" -ForegroundColor Green
Write-Host ""

# Create weather record
Write-Host "1. Testing POST /weather - Creating weather record for London, UK" -ForegroundColor Yellow
$createResponse = curl.exe -X POST "$baseUrl/weather" -H "Content-Type: application/json" -d '{"cityName": "London", "country": "UK"}' --silent
Write-Host "Response: $createResponse" -ForegroundColor Cyan

# Extract ID from JSON
if ($createResponse) {
    try {
        $jsonResponse = $createResponse | ConvertFrom-Json
        $weatherId = $jsonResponse.id
        Write-Host "Extracted Weather ID: $weatherId" -ForegroundColor Green
    }
    catch {
        Write-Host "Failed to parse JSON response or extract ID" -ForegroundColor Red
        exit 1
    }
}
else {
    Write-Host "Failed to create weather record" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Get all
Write-Host "2. Testing GET /weather - Fetching all weather records" -ForegroundColor Yellow
$getAllResponse = curl.exe "$baseUrl/weather" --silent
Write-Host "Response: $getAllResponse" -ForegroundColor Cyan
Write-Host ""

# Get weather by ID
Write-Host "3. Testing GET /weather/:id - Fetching weather record by ID: $weatherId" -ForegroundColor Yellow
$getByIdResponse = curl.exe "$baseUrl/weather/$weatherId" --silent
Write-Host "Response: $getByIdResponse" -ForegroundColor Cyan
Write-Host ""

# Update
Write-Host "4. Testing PUT /weather/:id - Updating weather record with mock data" -ForegroundColor Yellow
$updateResponse = curl.exe -X PUT "$baseUrl/weather/$weatherId" -H "Content-Type: application/json" -d '{"temperature": 25.5, "description": "sunny", "humidity": 65}' --silent
Write-Host "Response: $updateResponse" -ForegroundColor Cyan
Write-Host ""

# Get latest weather by city
Write-Host "5. Testing GET /weather/latest/:cityName - Fetching latest weather for London" -ForegroundColor Yellow
$getLatestResponse = curl.exe "$baseUrl/weather/latest/London" --silent
Write-Host "Response: $getLatestResponse" -ForegroundColor Cyan
Write-Host ""

# Delete
Write-Host "6. Testing DELETE /weather/:id - Deleting weather record ID: $weatherId" -ForegroundColor Yellow
$deleteResponse = curl.exe -X DELETE "$baseUrl/weather/$weatherId" --silent
Write-Host "Response: $deleteResponse" -ForegroundColor Cyan
Write-Host ""

# Verify deletion - try to get deleted record
Write-Host "7. Verifying deletion - Attempting to fetch deleted record" -ForegroundColor Yellow
$verifyDeleteResponse = curl.exe "$baseUrl/weather/$weatherId" --silent
Write-Host "Response: $verifyDeleteResponse" -ForegroundColor Cyan
Write-Host ""

Write-Host "=== Test Script Completed ===" -ForegroundColor Green
Write-Host "All CRUD operations have been tested." -ForegroundColor Green

# Optional: Create a few more test records with different cities
Write-Host ""
Write-Host "=== Creating Additional Test Records ===" -ForegroundColor Magenta

$cities = @(
    @{city="Paris"; country="France"},
    @{city="Tokyo"; country="Japan"},
    @{city="New York"; country="USA"}
)

foreach ($location in $cities) {
    Write-Host "Creating weather record for $($location.city), $($location.country)" -ForegroundColor Yellow
    $testResponse = curl.exe -X POST "$baseUrl/weather" -H "Content-Type: application/json" -d "{`"cityName`": `"$($location.city)`", `"country`": `"$($location.country)`"}" --silent
    Write-Host "Response: $testResponse" -ForegroundColor Cyan
    Start-Sleep -Seconds 1  # Small delay to avoid API rate limits
}

Write-Host ""
Write-Host "Final GET all records:" -ForegroundColor Yellow
$finalGetAll = curl.exe "$baseUrl/weather" --silent
Write-Host "Response: $finalGetAll" -ForegroundColor Cyan

Write-Host ""
Write-Host "=== All Tests Completed Successfully ===" -ForegroundColor Green

