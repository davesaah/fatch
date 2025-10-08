package config

import (
	_ "embed"
	"encoding/json"
	"errors"
)

//go:embed database.json
var databaseJSON []byte

type DB struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	DBName   string `json:"dbname"`
	Port     int    `json:"port"`
}

func LoadDBConfig() (*DB, error) {
	var db DB
	if err := json.Unmarshal(databaseJSON, &db); err != nil {
		return nil, errors.New("Unable to load database config")
	}

	return &db, nil
}
