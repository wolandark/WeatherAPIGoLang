package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"WeatherAPITask/internal/models"
	"WeatherAPITask/pkg/utils"
)

type WeatherService struct {
	db *gorm.DB
}

func NewWeatherService(db *gorm.DB) *WeatherService {
	return &WeatherService{db: db}
}

func (s *WeatherService) GetAll() ([]models.Weather, error) {
	var weather []models.Weather
	result := s.db.Find(&weather)
	return weather, result.Error
}

func (s *WeatherService) GetByID(id string) (*models.Weather, error) {
	var weather models.Weather
	result := s.db.First(&weather, "id = ?", id)
	return &weather, result.Error
}

func (s *WeatherService) GetLatestByCity(cityName string) (*models.Weather, error) {
	var weather models.Weather
	result := s.db.Where("city_name = ?", cityName).Order("fetched_at desc").First(&weather)
	return &weather, result.Error
}

func (s *WeatherService) Create(req models.WeatherRequest) (*models.Weather, error) {
	weatherData, err := s.fetchWeatherFromAPI(req.CityName, req.Country)
	if err != nil {
		return nil, err
	}

	weather := models.Weather{
		ID:          uuid.New().String(),
		CityName:    weatherData.Location.Name,
		Country:     weatherData.Location.Country,
		Temperature: weatherData.Current.TempC,
		Description: weatherData.Current.Condition.Text,
		Humidity:    weatherData.Current.Humidity,
		WindSpeed:   weatherData.Current.WindKph,
		FetchedAt:   time.Now(),
	}

	result := s.db.Create(&weather)
	return &weather, result.Error
}

func (s *WeatherService) Update(id string, updateData map[string]interface{}) (*models.Weather, error) {
	var weather models.Weather
	result := s.db.First(&weather, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	if temp, ok := updateData["temperature"]; ok {
		if tempFloat, err := utils.ParseFloat(temp); err == nil {
			weather.Temperature = tempFloat
		}
	}
	if desc, ok := updateData["description"]; ok {
		if descStr, ok := desc.(string); ok {
			weather.Description = descStr
		}
	}
	if hum, ok := updateData["humidity"]; ok {
		if humInt, err := utils.ParseInt(hum); err == nil {
			weather.Humidity = humInt
		}
	}
	if wind, ok := updateData["windSpeed"]; ok {
		if windFloat, err := utils.ParseFloat(wind); err == nil {
			weather.WindSpeed = windFloat
		}
	}

	result = s.db.Save(&weather)
	return &weather, result.Error
}

func (s *WeatherService) Delete(id string) error {
	result := s.db.Delete(&models.Weather{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (s *WeatherService) fetchWeatherFromAPI(city, country string) (*models.WeatherAPIResponse, error) {
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

	var weatherData models.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}
