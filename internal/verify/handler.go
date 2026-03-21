package verify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SendPayload struct {
	Email string `json:"email"`
}

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	var payload SendPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Некорректный payload", http.StatusBadRequest)
		return
	}
	hash := h.Service.GenerateHash()

	if err := h.Service.SaveHash(hash, payload.Email); err != nil {
		http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
		return
	}

	link := fmt.Sprintf("http://localhost:8081/verify/%s", hash)
	if err := h.Service.SendEmail(payload.Email, link); err != nil {
		http.Error(w, "Не удалось отправить письмо", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Письмо отправлено"))
}

func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	prefix := "/verify/"
	if len(path) <= len(prefix) {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}
	id := path[len(prefix):]
	email, ok := h.Service.VerifyHash(id)
	if !ok {
		http.Error(w, "false", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "true: %s подтвержден", email)
}
