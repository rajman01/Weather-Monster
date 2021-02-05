package routes

import (
	"fmt"
	"net/http"
	"strconv"
	m "weather_monster/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func CreateTemperature(c *gin.Context) {
	temperature := m.Temperature{}
	err := c.ShouldBindJSON(&temperature)
	fmt.Println(temperature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	_, erro := m.FindCityById(temperature.City_ID, &conn)

	if erro != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": erro.Error()})
		return
	}

	err = temperature.AddTemperature(&conn)
	if err != nil {
		fmt.Println("Error in .AddTemperature()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, temperature)
}

func ForcastRoute(c *gin.Context) {
	cityID, _ := strconv.ParseInt(c.Param("city_Id"), 10, 64)
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	forcast := m.Forcast{}
	_, err := m.FindCityById(cityID, &conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	forcast.City_ID = cityID

	temperatures, err := forcast.GetForcast(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var max float64
	var min float64
	for _, e := range temperatures {
		max += float64(e.Max)
		min += float64(e.Min)
	}
	fmt.Println(max, min, len(temperatures))
	
	forcast.Max = float64(max) / float64(len(temperatures))
	forcast.Min = float64(min) / float64(len(temperatures))
	forcast.Sample = len(temperatures)

	c.JSON(http.StatusOK, forcast)
}
