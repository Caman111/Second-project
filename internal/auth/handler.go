package auth

import (
	"3-validation-api/internal/models"
	"fmt"
	"net/http"
)

func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(models.UserPhoneKey).(string)
	if !ok {
		http.Error(w, "Внутренняя ошибка: не удалось получить пользователя", http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf("Пользователь %s успешно совершил покупку!", userID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	phone, ok := r.Context().Value(models.UserPhoneKey).(string)
	if !ok {
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "success", "phone": "%s"}`, phone)
}
