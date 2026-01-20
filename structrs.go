package main

import "time"

type Chat struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"size:200;not null"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	ChatID    uint      `gorm:"not null; index"`
	Text      string    `gorm:"size:5000; not null`
	CreatedAt time.Time `gorm:"default:now()"`
	Chat      Chat      `gorm:"constraint:OnDelete:CASCADE;foreignKey:ChatID"`
}

type MessageText struct {
	Text string `json:"text"`
}
