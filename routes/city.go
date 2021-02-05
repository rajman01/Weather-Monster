package routes

import (
	"fmt"
	"net/http"
	"strconv"
	m "weather_monster/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func CreateCity(c *gin.Context) {
	city := m.City{}
	err := c.ShouldBindJSON(&city)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = city.AddCity(&conn)
	if err != nil {
		fmt.Println("Error in city.Addcity()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, city)
}

func CityUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("Id"), 10, 64)
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	citySent := m.City{}
	err := c.ShouldBindJSON(&citySent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form sent"})
		return
	}

	cityBeingUpdated, err := m.FindCityById(id, &conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	citySent.ID = cityBeingUpdated.ID
	err = citySent.UpdateCity(&conn, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, citySent)
}

func CityDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("Id"), 10, 64)
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	citySent := m.City{}
	cityBeingDeleted, err := m.FindCityById(id, &conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	citySent.ID = cityBeingDeleted.ID
	citySent.Name = cityBeingDeleted.Name
	citySent.Latitude = cityBeingDeleted.Latitude
	citySent.Longitude = cityBeingDeleted.Longitude

	err = citySent.DeleteCity(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, citySent)
}
