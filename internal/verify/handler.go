package verify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	h.Service.SaveHash(hash, payload.Email)

	link := fmt.Sprintf("http://localhost:8080/verify/%s", hash)
	if err := h.Service.SendEmail(payload.Email, link); err != nil {
		http.Error(w, "Не удалось отправить письмо", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Письмо отправлено"))

	data := map[string]string{payload.Email: hash}
	file, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile("users.json", file, 0644)

	fmt.Println("Данные сохранены в users.json")
}

func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	prefix := "/verify/"

	if len(path) <= len(prefix) {
		http.Error(w, "Неверный URL, нужно /verify/ID", http.StatusBadRequest)
		return
	}

	id := path[len(prefix):]
	fmt.Printf("Пришел запрос на верификацию. ID: %s\n", id)
	fmt.Fprintf(w, "Проверка для ID: %s\n", id)

}
