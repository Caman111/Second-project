package handler

import (
	"3-validation-api/internal/models"
	"net/http"
)

type OrderHandler struct{}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	phone, ok := r.Context().Value(models.UserPhoneKey).(string)
	if !ok {
		http.Error(w, "User phone not found", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Заказ оформлен на номер: " + phone))
}
