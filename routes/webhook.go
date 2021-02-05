package routes

import (
	"fmt"
	"net/http"
	"strconv"
	m "weather_monster/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func CreateWebhook(c *gin.Context) {
	webhook := m.Webhook{}
	err := c.ShouldBindJSON(&webhook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	_, erro := m.FindCityById(webhook.City_ID, &conn)

	if erro != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": erro.Error()})
		return
	}

	err = webhook.AddWebhook(&conn)
	if err != nil {
		fmt.Println("Error in .AddWebhook()")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, webhook)
}

func WebhookDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("Id"), 10, 64)
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	webhookSent := m.Webhook{}
	webhookBeingDeleted, err := m.FindWebhookById(id, &conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhookSent.ID = webhookBeingDeleted.ID
	webhookSent.Callback_URL = webhookBeingDeleted.Callback_URL
	webhookSent.City_ID = webhookBeingDeleted.City_ID
	
	err = webhookSent.DeleteWebhook(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhookSent)
}
