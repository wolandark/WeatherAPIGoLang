# Weather API

A REST API for storing and retrieving weather data using Go, Gin, and GORM with Mariadb database.

## Prerequisites

- Go 24.04
- Mariadb
- WeatherAPI.com account and API key

## Setup

1. Clone the repository and navigate to the project directory

2. Install dependencies:
```bash
go mod tidy
```

3. Create a Mariadb database using docker:
```bash
docker run -d --name mariadb-weather -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=weather_db -p 3306:3306 mariadb:latest
```

4. Configure environment variables by creating a `.env` file:
```env
WEATHER_API_KEY=apikey
DB_USER=user
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=weather_db

PORT=8080
```

A .env.example file is provided as a guide.

5. Run the application:
```bash
go run main.go
```

The server will start on the specified port (default: 8080) and automatically create the required database tables.

## API Endpoints

### GET /weather
Returns all weather records.

### GET /weather/:id
Returns a specific weather record by ID.

### GET /weather/latest/:cityName
Returns the most recent weather record for a city.

### POST /weather
Creates a new weather record by fetching current data from WeatherAPI.

Request body:
```json
{
  "cityName": "London",
  "country": "UK"
}
```

### PUT /weather/:id
Updates an existing weather record.

Request body (partial updates supported):
```json
{
  "temperature": 25.5,
  "description": "Sunny",
  "humidity": 60,
  "windSpeed": 10.2
}
```

### DELETE /weather/:id
Deletes a weather record by ID.

## Testing

Automated test scripts are available in the `Tests` folder:
- Linux/macOS: Run the shell script
    - The Linux Bash script requires the jq command line tool
- Windows: Run the powershell script
    - The powershell test script must be run from powershell (duh). Its best to use powershell 7

Both scripts require cURL. These scripts execute curl commands to test all API endpoints.

