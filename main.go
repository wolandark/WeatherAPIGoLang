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

type OpenWeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Name string `json:"name"`
}











func main() {
	//
}
