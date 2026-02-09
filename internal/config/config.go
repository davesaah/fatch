// package config loads the .env.json file for use in the API
package config

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	User     string `json:"USER"`
	Password string `json:"PASSWORD"`
	Host     string `json:"HOST"`
	Port     string `json:"PORT"`
	Name     string `json:"NAME"`
	Schema   string `json:"SCHEMA"`
}

type Config struct {
	DBConfig    DBConfig `json:"DB"`
	Environment string   `json:"ENVIRONMENT"`
	JWTSecret   string   `json:"JWT_SECRET"`
}

// LoadConfig loads the .env file and returns a Config struct
func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(".env.json")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
