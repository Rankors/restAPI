package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func create_chat(w http.ResponseWriter, r *http.Request) {
	var chat Chat

	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "ошибка при декоде Title", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(chat.Title) == "" {
		http.Error(w, "пустой Title", http.StatusBadRequest)
		return
	}
	res := db.Create(&chat)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("чат создан, id:%d", chat.ID)
	sendJSON(w, http.StatusOK, map[string]uint{"id": chat.ID})
}

func send_message(w http.ResponseWriter, r *http.Request, ChatID string) {

	id, err := strconv.Atoi(ChatID)
	if err != nil {
		http.Error(w, "это не id", http.StatusBadRequest)
		return
	}
	chat, err := find_chat(id)
	if err != nil {
		http.Error(w, "не существующий ID", http.StatusNotFound)
		return
	}
	var txt MessageText
	if err := json.NewDecoder(r.Body).Decode(&txt); err != nil || strings.TrimSpace(txt.Text) == "" {
		http.Error(w, "Пустой text", http.StatusBadRequest)
		return
	}
	msg := Message{
		ChatID: chat.ID,
		Text:   txt.Text,
		Chat:   *chat,
	}
	res := db.Create(&msg)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, http.StatusOK, msg)
	log.Println("Новое сообщение: %s", msg.Text)
}

func get_chat(w http.ResponseWriter, r *http.Request, ChatID string) {

	id, err := strconv.Atoi(ChatID)
	if err != nil {
		http.Error(w, "это не id", http.StatusBadRequest)
		return
	}
	chat, err := find_chat(id)
	if err != nil {
		http.Error(w, "Такой чат не найден", http.StatusNotFound)
		return
	}

	limit := 20
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		if l > 100 {
			limit = 100
		} else if l > 0 {
			limit = l
		}
	}

	var msg []Message

	if err := db.Where("chat_id = ?", id).Order("created_at DESC").Limit(limit).Find(&msg).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, map[string]any{
		"chat":     chat,
		"messages": msg,
	})
	log.Println("запросили чат :%d", id)
}

func delete_chat(w http.ResponseWriter, r *http.Request, ChatID string) {
	id, err := strconv.Atoi(ChatID)
	if err != nil {
		http.Error(w, "это не id", http.StatusBadRequest)
		return
	}
	chat, err := find_chat(id)
	if err != nil {
		http.Error(w, "Такой чат не найден", http.StatusNotFound)
		return
	}

	db.Delete(&chat, id)

	sendJSON(w, http.StatusNoContent, 0)
	log.Println("удалили чат :%d", id)
}
