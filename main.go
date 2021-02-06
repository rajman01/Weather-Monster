package main

import (
	"context"
	"fmt"
	r "weather_monster/routes"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func main() {
	conn, err := connectDB()
	if err != nil {
		return
	}
	router := gin.Default()
	router.Use(dbMiddleware(*conn))

	citiesGroup := router.Group("cities")
	{
		citiesGroup.POST("", r.CreateCity)
		citiesGroup.PATCH("/:Id", r.CityUpdate)
		citiesGroup.DELETE("/:Id", r.CityDelete)
	}

	tempraturesGroup := router.Group("temperatures")
	{
		tempraturesGroup.POST("", r.CreateTemperature)
	}

	forcastGroup := router.Group("forecasts")
	{
		forcastGroup.GET("/:city_Id", r.ForcastRoute)
	}

	webhooksGroup := router.Group("webhooks")
	{
		webhooksGroup.POST("", r.CreateWebhook)
		webhooksGroup.DELETE("/:Id", r.WebhookDelete)
	}

	router.Run(":3000")
}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgres://aqndmbyupzebhj:d51ea7fb306217099832c92802e158e307a5a724a88610e75d2cbf5c1e17a26b@ec2-18-204-101-137.compute-1.amazonaws.com:5432/d4hnjei0qe63db")
	if err != nil {
		fmt.Println("Error connecting to db")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}
