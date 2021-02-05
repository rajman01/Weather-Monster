package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Webhook struct {
	ID        int64   `json:"id"`
	City_ID   int64     `json:"city_id"`
	Callback_URL  string `json:"callback_url"`
}

func (w *Webhook) AddWebhook(conn *pgx.Conn) error {

	w.Callback_URL = strings.ToLower(strings.TrimSpace(w.Callback_URL))
	if len(w.Callback_URL) < 1 {
		return fmt.Errorf("url cannot be blank")
	}
	
	row := conn.QueryRow(context.Background(), "INSERT INTO webhook (city_id, callback_url) VALUES($1, $2) RETURNING id", w.City_ID, w.Callback_URL)
	err := row.Scan(&w.ID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating webhook")
	}
	return nil
}

func (w *Webhook) DeleteWebhook(conn *pgx.Conn) error {

	_, err := conn.Exec(context.Background(), "DELETE FROM webhook WHERE id=$1", w.ID)

	if err != nil {
		fmt.Printf("Error deleting webhook: (%v)", err)
		return fmt.Errorf("Error deleting webhook")
	}
	return nil
}


func FindWebhookById(id int64, conn *pgx.Conn) (Webhook, error) {
	row := conn.QueryRow(context.Background(), "SELECT city_id, callback_url FROM webhook WHERE id=$1", id)
	webhook := Webhook{
		ID: id,
	}
	err := row.Scan(&webhook.City_ID, &webhook.Callback_URL)
	if err != nil {
		return webhook, fmt.Errorf("The webhook doesn't exist")
	}

	return webhook, nil
}