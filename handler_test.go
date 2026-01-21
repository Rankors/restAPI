package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var CreatedChatID string

func TestMain(m *testing.M) {
	initDB() // инициализируем базу, без нее никуда
	db.AutoMigrate(&Chat{}, &Message{})
	code := m.Run()
	os.Exit(code)
}

func zapros(data string, method string, address string) (*httptest.ResponseRecorder, *http.Request) {
	body := bytes.NewBuffer([]byte(data))
	req := httptest.NewRequest(method, address, body)
	w := httptest.NewRecorder()
	return w, req
}

func TestCreateChat(t *testing.T) {
	w, req := zapros(`{"title":"name chat"}`, http.MethodPost, "/chats/")
	create_chat(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "должен вернуться 201 Created")
	assert.Contains(t, string(body), "name chat", "ответ должен содержать название чата")
	var res map[string]any
	_ = json.Unmarshal(body, &res)
	CreatedChatID = fmt.Sprint(res["id"])
}

func TestSendMessage(t *testing.T) {
	w, req := zapros(`{"text":"Привет"}`, http.MethodPost, "/chats/"+CreatedChatID+"/messages")
	send_message(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "должен вернуться 201 Created")
	assert.Contains(t, string(body), "Привет", "ответ должен содержать текст сообщения")
}

func TestDeleteChat(t *testing.T) {
	w, req := zapros(`{}`, http.MethodDelete, "/chats/"+CreatedChatID)
	delete_chat(w, req, "1")

	resp := w.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "должен вернуться 204 No Content")
}

func TestGetChat(t *testing.T) {
	w, req := zapros(`{}`, http.MethodGet, "/chats/"+CreatedChatID+"?limit=3")
	get_chat(w, req, "1")

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "должен вернуться 200 OK")
	assert.Contains(t, string(body), "messages", "ответ должен содержать список сообщений")
}
