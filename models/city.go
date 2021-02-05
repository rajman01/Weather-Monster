package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
)

type City struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *City) AddCity(conn *pgx.Conn) error {

	c.Name = strings.ToLower(strings.TrimSpace(c.Name))
	if len(c.Name) < 1 {
		return fmt.Errorf("Name cannot be blank")
	}
	row := conn.QueryRow(context.Background(), "SELECT id from city WHERE name = $1", c.Name)
	cityLookup := City{}
	err := row.Scan(&cityLookup)
	if err != pgx.ErrNoRows {
		fmt.Println("found city")
		return fmt.Errorf("A city with this name already exists")
	}
	row = conn.QueryRow(context.Background(), "INSERT INTO city (name, latitude, longitude) VALUES($1, $2, $3) RETURNING id", c.Name, c.Latitude, c.Longitude)
	erro := row.Scan(&c.ID)

	if erro != nil {
		fmt.Println(erro)
		return fmt.Errorf("There was an error creating the city")
	}
	return nil
}

func (c *City) UpdateCity(conn *pgx.Conn, Id int64) error {
	c.Name = strings.ToLower(strings.TrimSpace(c.Name))
	if len(c.Name) < 1 {
		return fmt.Errorf("Name cannot be blank")
	}
	row := conn.QueryRow(context.Background(), "SELECT id from city WHERE name = $1", c.Name)
	cityLookup := City{}
	err := row.Scan(&cityLookup.ID)
	if err != pgx.ErrNoRows {
		fmt.Println("found city")
		if cityLookup.ID != Id {
			return fmt.Errorf("A city with this name already exists")
		}
	}

	_, erro := conn.Exec(context.Background(), "UPDATE city SET name=$1, latitude=$2, longitude=$3 WHERE id=$4", c.Name, c.Latitude, c.Longitude, c.ID)

	if erro != nil {
		fmt.Printf("Error updating city: (%v)", erro)
		return fmt.Errorf("Error updating city")
	}
	return nil
}

func (c *City) DeleteCity(conn *pgx.Conn) error {

	_, err := conn.Exec(context.Background(), "DELETE FROM city WHERE id=$1", c.ID)

	if err != nil {
		fmt.Printf("Error deleting city: (%v)", err)
		return fmt.Errorf("Error deleting city")
	}
	return nil
}

func FindCityById(id int64, conn *pgx.Conn) (City, error) {
	row := conn.QueryRow(context.Background(), "SELECT name, latitude, longitude FROM city WHERE id=$1", id)
	city := City{
		ID: id,
	}
	err := row.Scan(&city.Name, &city.Latitude, &city.Longitude)
	if err != nil {
		return city, fmt.Errorf("The city doesn't exist")
	}

	return city, nil
}
