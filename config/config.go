// Package config loads any config data that the api requires
package config

import (
	_ "embed"
	"encoding/json"
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
		return nil, err
	}

	return &db, nil
}

//go:embed jwt.json
var jwtJSON []byte

type JWT struct {
	Secret string `json:"secret"`
}

func LoadJWTConfig() ([]byte, error) {
	var jwt JWT
	if err := json.Unmarshal(jwtJSON, &jwt); err != nil {
		return nil, err
	}

	secret := []byte(jwt.Secret)

	return secret, nil
}
