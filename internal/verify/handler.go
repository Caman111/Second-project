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
	var data struct {
		SessionID string `json:"sessionId"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	email, ok := h.Service.VerifyHash(data.SessionID)
	if !ok {
		http.Error(w, "false", 401)
		return
	}
	fmt.Fprintf(w, "true: %s", email)
}
