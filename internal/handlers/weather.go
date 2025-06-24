package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	
	"WeatherAPITask/internal/models"
	"WeatherAPITask/internal/services"
)

type WeatherHandler struct {
	weatherService *services.WeatherService
}

func NewWeatherHandler(db *gorm.DB) *WeatherHandler {
	return &WeatherHandler{
		weatherService: services.NewWeatherService(db),
	}
}

func (h *WeatherHandler) GetAllWeather(c *gin.Context) {
	weather, err := h.weatherService.GetAll()
	if err != nil {
		log.Printf("Error fetching weather records: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather records"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func (h *WeatherHandler) GetWeatherByID(c *gin.Context) {
	id := c.Param("id")
	
	weather, err := h.weatherService.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
			return
		}
		log.Printf("Error fetching weather record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather record"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func (h *WeatherHandler) GetLatestWeatherByCity(c *gin.Context) {
	cityName := c.Param("cityName")
	
	weather, err := h.weatherService.GetLatestByCity(cityName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No weather records found for this city"})
			return
		}
		log.Printf("Error fetching latest weather record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather record"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func (h *WeatherHandler) CreateWeather(c *gin.Context) {
	var req models.WeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weather, err := h.weatherService.Create(req)
	if err != nil {
		log.Printf("Error creating weather record: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create weather record"})
		return
	}

	log.Printf("Weather record created for %s, %s", weather.CityName, weather.Country)
	c.JSON(http.StatusCreated, weather)
}

func (h *WeatherHandler) UpdateWeather(c *gin.Context) {
	id := c.Param("id")
	
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weather, err := h.weatherService.Update(id, updateData)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
			return
		}
		log.Printf("Error updating weather record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weather record"})
		return
	}

	log.Printf("Weather record updated for ID: %s", id)
	c.JSON(http.StatusOK, weather)
}

func (h *WeatherHandler) DeleteWeather(c *gin.Context) {
	id := c.Param("id")

	err := h.weatherService.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Weather record not found"})
			return
		}
		log.Printf("Error deleting weather record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete weather record"})
		return
	}

	log.Printf("Weather record deleted for ID: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Weather record deleted successfully"})
}
