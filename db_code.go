package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() {
	dsn := getDSN()
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Chat{}, &Message{})
	log.Println("Подключение к PostgreSQL удачно")
}

func find_chat(id int) (*Chat, error) {
	var chat Chat
	if err := db.First(&chat, id).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func getDSN() string {
	data, _ := os.ReadFile("db.json")
	var cfg Config
	_ = json.Unmarshal(data, &cfg)

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
}
