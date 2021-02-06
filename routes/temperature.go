package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	webhooks, erroo := temperature.AddTemperature(&conn)
	if erroo != nil {
		fmt.Println("Error in temperature.AddTemperature()")
		c.JSON(http.StatusBadRequest, gin.H{"error": erroo.Error()})
		return
	}
	for _, e := range webhooks {
		jsonData := map[string]string{
			"city_id":   fmt.Sprint(temperature.City_ID),
			"min":       fmt.Sprint(temperature.Min),
			"max":       fmt.Sprint(temperature.Max),
			"timestamp": fmt.Sprint(temperature.Timestamp),
		}
		jsonValue, _ := json.Marshal(jsonData)
		fmt.Println(jsonValue)
		response, err2 := http.Post(e.Callback_URL, "application/json", bytes.NewBuffer(jsonValue))
		if err2 != nil {
			fmt.Printf("This Http request failed with error%s\n", err2)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
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
	fmt.Println(temperatures)
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
	if len(temperatures) != 0 {
		forcast.Max = float64(max) / float64(len(temperatures))
		forcast.Min = float64(min) / float64(len(temperatures))
	} else {
		forcast.Max = float64(0)
		forcast.Min = float64(0)
	}
	forcast.Sample = len(temperatures)

	c.JSON(http.StatusOK, forcast)
}
