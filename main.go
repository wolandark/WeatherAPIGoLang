package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Weather struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CityName    string    `json:"cityName" gorm:"not null"`
	Country     string    `json:"country" gorm:"not null"`
	Temperature float64   `json:"temperature" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Humidity    int       `json:"humidity" gorm:"not null"`
	WindSpeed   float64   `json:"windSpeed" gorm:"not null"`
	FetchedAt   time.Time `json:"fetchedAt" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type WeatherRequest struct {
	CityName string `json:"cityName" binding:"required"`
	Country  string `json:"country" binding:"required"`
}

type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		WindKph   float64 `json:"wind_kph"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" {
		dbUser = "root"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "weather_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	dbUser, dbPassword, dbHost, dbPort, dbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db = database
	db.AutoMigrate(&Weather{})

	log.Println("Database connected and migrated successfully")


	r := gin.Default()

	r.GET("/weather", getAllWeather)
	r.GET("/weather/:id", getWeatherByID)
	r.GET("/weather/latest/:cityName", getLatestWeatherByCity)

	r.POST("/weather", createWeather)
	r.PUT("/weather/:id", updateWeather)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}


func getAllWeather(c *gin.Context) {
	var weather []Weather
	result := db.Find(&weather)
	if result.Error != nil {
		log.Printf("Error fetching weather records: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather records"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func getWeatherByID(c *gin.Context) {
	id := c.Param("id")
	var weather Weather

	result := db.First(&weather, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
			return
		}
		log.Printf("Error fetching weather record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather record"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func getLatestWeatherByCity(c *gin.Context) {
	cityName := c.Param("cityName")
	var weather Weather

	result := db.Where("city_name = ?", cityName).Order("fetched_at desc").First(&weather)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No weather records found for this city"})
			return
		}
		log.Printf("Error fetching latest weather record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather record"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func createWeather(c *gin.Context) {
	var req WeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weatherData, err := fetchWeatherFromAPI(req.CityName, req.Country)
	if err != nil {
		log.Printf("Error fetching weather from API: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch weather data"})
		return
	}

	weather := Weather{
		ID:          uuid.New().String(),
		CityName:    weatherData.Location.Name,
		Country:     weatherData.Location.Country,
		Temperature: weatherData.Current.TempC,
		Description: weatherData.Current.Condition.Text,
		Humidity:    weatherData.Current.Humidity,
		WindSpeed:   weatherData.Current.WindKph,
		FetchedAt:   time.Now(),
	}

	result := db.Create(&weather)
	if result.Error != nil {
		log.Printf("Error saving weather record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save weather record"})
		return
	}

	log.Printf("Weather record created for %s, %s", weather.CityName, weather.Country)
	c.JSON(http.StatusCreated, weather)
}

func fetchWeatherFromAPI(city, country string) (*WeatherAPIResponse, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s", apiKey, city, country)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WeatherAPI returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherData WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}

func updateWeather(c *gin.Context) {
	id := c.Param("id")
	var weather Weather

	result := db.First(&weather, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
			return
		}
		log.Printf("Error fetching weather record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather record"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if temp, ok := updateData["temperature"]; ok {
		if tempFloat, err := parseFloat(temp); err == nil {
			weather.Temperature = tempFloat
		}
	}
	if desc, ok := updateData["description"]; ok {
		if descStr, ok := desc.(string); ok {
			weather.Description = descStr
		}
	}
	if hum, ok := updateData["humidity"]; ok {
		if humInt, err := parseInt(hum); err == nil {
			weather.Humidity = humInt
		}
	}
	if wind, ok := updateData["windSpeed"]; ok {
		if windFloat, err := parseFloat(wind); err == nil {
			weather.WindSpeed = windFloat
		}
	}

	result = db.Save(&weather)
	if result.Error != nil {
		log.Printf("Error updating weather record: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weather record"})
		return
	}

	log.Printf("Weather record updated for ID: %s", id)
	c.JSON(http.StatusOK, weather)
}

func parseFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert to float64")
	}
}

func parseInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert to int")
	}
}
