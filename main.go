package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

var server *http.Server
var db *gorm.DB

func main() {
	initDB()
	http.HandleFunc("/", router)
	fmt.Println("Server started at :8080")
	log.Println("Сервер слушает")
	http.ListenAndServe(":8080", nil)
}

func router(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// удаляем в конце / разбиваем на части, разделитель /
	switch len(path) {
	case 1: // считаем по длинне например /chats/ длинна 1
		if r.Method == http.MethodPost {
			create_chat(w, r)
		}
	case 2: // длинна 2 уже chats/id
		if r.Method == http.MethodGet {
			get_chat(w, r, path[1])
		}
		if r.Method == http.MethodDelete {
			delete_chat(w, r, path[1])
		}
	case 3: // и тд
		if r.Method == http.MethodPost {
			send_message(w, r, path[1])
		}
	default: // обращение не по тому адрессу
		fmt.Fprintln(w, "Нет такого адресса запроса")
	}
}

func sendJSON(w http.ResponseWriter, status int, data any) { // отправка готовый жсон клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
