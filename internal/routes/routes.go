package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	
	"WeatherAPITask/internal/handlers"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	weatherHandler := handlers.NewWeatherHandler(db)

	// Weather routes
	r.GET("/weather", weatherHandler.GetAllWeather)
	r.GET("/weather/:id", weatherHandler.GetWeatherByID)
	r.GET("/weather/latest/:cityName", weatherHandler.GetLatestWeatherByCity)
	r.POST("/weather", weatherHandler.CreateWeather)
	r.PUT("/weather/:id", weatherHandler.UpdateWeather)
	r.DELETE("/weather/:id", weatherHandler.DeleteWeather)
}
