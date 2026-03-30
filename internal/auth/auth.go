package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Repo *AuthRepository
}

type LoginRequest struct {
	Phone string `json:"phone"`
}

type LoginResponse struct {
	SessionID string `json:"sessionId"`
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Доступ разрешен."}`))
}

func sendRealSMS(phone, code string) error {
	apiID := os.Getenv("SMS_API_ID")
	if apiID == "" {
		return fmt.Errorf("SMS_API_ID не найден в .env")
	}

	msg := url.PathEscape(code)

	apiURL := fmt.Sprintf("https://sms.ru/sms/send?api_id=%s&to=%s&msg=%s&json=1",
		apiID, phone, msg)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("ошибка сети: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("ошибка парсинга: %v", err)
	}

	fmt.Printf("\n--- DEBUG SMS.RU (SMS MODE) ---\nОтвет: %+v\n--------------------\n\n", result)

	return nil
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Ошибка в формате JSON", http.StatusBadRequest)
			return
		}

		if err := h.Repo.CheckLimit(req.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		sessionID := uuid.New().String()
		code := fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
		fmt.Printf("\n[DEV MODE] КОД ДЛЯ ВХОДА: %s\n\n", code)

		if err := sendRealSMS(req.Phone, code); err != nil {
			fmt.Println("Ошибка SMS.ru:", err)
			http.Error(w, "Не удалось отправить SMS", http.StatusInternalServerError)
			return
		}

		h.Repo.SaveSession(sessionID, Session{
			Phone: req.Phone,
			Code:  code,
		})
		h.Repo.UpdateLimit(req.Phone)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{SessionID: sessionID})
	}
}

func (h *AuthHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SessionID string `json:"sessionId"`
			Code      string `json:"code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": false, "message": "Invalid JSON"})
			return
		}
		session, ok := h.Repo.GetSession(req.SessionID)
		fmt.Println("\n[DEBUG VERIFY START]")
		fmt.Printf("Ищем SessionID: %s\n", req.SessionID)
		if !ok {
			fmt.Println("РЕЗУЛЬТАТ: Сессия НЕ НАЙДЕНА в базе. Возможно, рестарт сервера стер мапу?")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "false", "message": "session not found"})
			return
		}

		fmt.Printf("В базе код: [%s] | Ты прислал: [%s]\n", session.Code, req.Code)

		// 3. Проверяем код
		if session.Code != req.Code {
			fmt.Println("РЕЗУЛЬТАТ: Коды не совпали!")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "false", "message": "wrong code"})
			return
		}

		fmt.Println("РЕЗУЛЬТАТ: Успех, генерирую JWT...")
		secretKey := []byte("JWT_SECRET")
		claims := jwt.MapClaims{
			"phone": session.Phone,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"iat":   time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		finalToken, err := token.SignedString(secretKey)
		if err != nil {
			http.Error(w, "Error signing token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"token":  finalToken,
		})
	}
}
