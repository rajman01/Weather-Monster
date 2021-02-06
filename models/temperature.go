package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Temperature struct {
	ID        int64     `json:"id"`
	City_ID   int64     `json:"city_id"`
	Max       float64   `json:"max"`
	Min       float64   `json:"min"`
	Timestamp time.Time `json:"timestamp"`
}

type Forcast struct {
	City_ID int64   `json:"city_id"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Sample  int     `json:"sample"`
}

func (t *Temperature) AddTemperature(conn *pgx.Conn) ([]Webhook, error) {

	row := conn.QueryRow(context.Background(), "INSERT INTO temperature (city_id, max, min) VALUES($1, $2, $3) RETURNING id, timestamp", t.City_ID, float64(t.Max), float64(t.Min))
	err := row.Scan(&t.ID, &t.Timestamp)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("There was an error creating the temperature")
	}

	rows, err := conn.Query(context.Background(), "SELECT * FROM webhook WHERE city_id =$1", t.City_ID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Something Happened")
	}

	var webhooks []Webhook
	for rows.Next() {
		webhook := Webhook{}
		err = rows.Scan(&webhook.ID, &webhook.City_ID, &webhook.Callback_URL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

func (f *Forcast) GetForcast(conn *pgx.Conn) ([]Temperature, error) {
	rows, err := conn.Query(context.Background(), "SELECT * FROM temperature WHERE city_id =$1 AND timestamp > NOW() - INTERVAL '24 hours'", f.City_ID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting forecast")
	}

	var temperatures []Temperature
	for rows.Next() {
		temperature := Temperature{}
		err = rows.Scan(&temperature.ID, &temperature.City_ID, &temperature.Max, &temperature.Min, &temperature.Timestamp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		temperatures = append(temperatures, temperature)
	}
	return temperatures, nil
}
