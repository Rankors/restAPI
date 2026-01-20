package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() {
	dsn := "host=localhost user=user password=password dbname=chatsdb port=5432 sslmode=disable"
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
