// Contains configuration files and code related to application configuration
package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type ServerConfig struct {
	Port       string           `json:"SERVER_PORT"`
	PostgresDB PostgresDBconfig `json:"POSTGRES_DB"`
	JwtSecret  string           `json:"JWT_SECRET"`
}

var serverConfigInstance *ServerConfig
var once sync.Once

type PostgresDBconfig struct {
	HOST     string `json:"host"`
	PORT     int    `json:"port"`
	USER     string `json:"user"`
	PASSWORD string `json:"password"`
	DBNAME   string `json:"dbname"`
}

func GetInstance() *ServerConfig {
	once.Do(func() {
		serverConfigInstance = &ServerConfig{}
	})
	return serverConfigInstance
}

func (s *ServerConfig) InitServerConfig() {

	file, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal("Error in reading config file ", err)
	}

	if err := json.Unmarshal(file, &s); err != nil {
		log.Fatal("Error parsing config file:", err)
	}
}

func (s *ServerConfig) GetServerPort() string {
	if s.Port != "" {
		return s.Port
	}
	return "8082"
}

func (s *ServerConfig) GetPostgresDBConf() *PostgresDBconfig {
	return &s.PostgresDB
}

func (s *ServerConfig) GetSecretKey() string {
	return s.JwtSecret
}
