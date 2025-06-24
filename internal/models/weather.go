package models

import (
	"time"
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
