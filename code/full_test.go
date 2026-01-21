package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullApi(t *testing.T) { // комплексная проверка всего API
	w, r := zapros(`{"title":"full test"}`, http.MethodPost, "/chats/")
	create_chat(w, r)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	var chat map[string]any
	_ = json.Unmarshal(body, &chat)
	chatId := fmt.Sprint(chat["id"]) // создали новый чат и получили его id

	for i := 1; i <= 10; i++ {
		payload := fmt.Sprintf(`{"text":"msg %d"}`, i)
		w, req := zapros(payload, http.MethodPost, "/chats/"+chatId+"/messages")
		send_message(w, req, chatId)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var m Message
		_ = json.Unmarshal(body, &m)
		assert.Equal(t, fmt.Sprintf("msg %d", i), m.Text)
		assert.Equal(t, chatId, fmt.Sprint(m.ChatID))
	} // отправляет сообщение по циклу

	w, r = zapros(`{}`, http.MethodGet, "/chats/"+chatId+"?limit=3") // запрос сообщений
	get_chat(w, r, chatId)
	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)
	var res map[string]any
	_ = json.Unmarshal(body, &res)

	messages := res["messages"].([]any)
	assert.Equal(t, 3, len(messages), "должно быть 3 сообщения") // полученные сообщения, если получим

	w, r = zapros(`{}`, http.MethodDelete, "/chats/"+chatId) // пробуем удалить чат
	delete_chat(w, r, chatId)
	resp = w.Result()
	assert.Equal(t, 204, resp.StatusCode)
}
