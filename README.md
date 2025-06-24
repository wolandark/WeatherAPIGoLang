# Weather API

A REST API for storing and retrieving weather data using Go, Gin, and GORM with Mariadb database.

# Notes For The Reviewer

MariaDB is preferred over PostgreSQL due to the simple fact that I already had its Go module installed. JWT authentication, though rudimentary and fairly straightforward to implement, was skipped due to current internet disruptions. The same applies to deployment.

No AI agents, LLMs, or similar tools were used in coding this project, except for translating my test bash script to PowerShell, which I could have done manually, but who has time for PowerShell, right?

I also had to choose https://www.weatherapi.com/ over OpenWeatherMap due to the same internet issues. It's free, easy to sign up for, and does the job perfectly for this demo.

Dockerization would have been included if more time was available, yet war has delayed us enough. A docker-compose setup is a manageable task, I'll provide one, though I cannot thoroughly test it. I'm fairly confident it will work.

AI was used to fix minor syntax issues in docker-compose.yml, since I currently have no yml linter installed. (Even if I did , I take no shame in using ai for yml, yml is an insult to humanity itself)

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
go run cmd/server/main.go
```

# Running with Docker

Use `docker compose build` to build the image. 

Use `docker compose up` to run the image.

Or:

`docker-compose up --build`

Use cURL or postman to test. 

Ex: `http://localhost:8080/weather`

The Provided test scripts should also work.

The `.env` file may also need to be updated when using the docker compose.
```env
WEATHER_API_KEY=key  
DB_USER=root
DB_PASSWORD=rootpassword
DB_HOST=db
DB_PORT=3306
DB_NAME=weather_db
PORT=8080
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

