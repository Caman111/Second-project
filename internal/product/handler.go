package product

import (
	models "3-validation-api/internal/models"
	"3-validation-api/pkg/res"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.Product

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res.Json(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	v := validator.New()
	if err := v.Struct(input); err != nil {
		res.Json(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&input).Error; err != nil {
		res.Json(w, http.StatusInternalServerError, map[string]string{"error": "DB error"})
		return
	}
	res.Json(w, http.StatusCreated, input)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p models.Product

	if err := h.DB.First(&p, id).Error; err != nil {
		res.Json(w, http.StatusNotFound, map[string]string{"error": "Товар не найден"})
		return
	}

	res.Json(w, http.StatusOK, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.Product

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res.Json(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	var p models.Product
	if err := h.DB.First(&p, id).Error; err != nil {
		res.Json(w, http.StatusNotFound, map[string]string{"error": "Товар не найден"})
		return
	}
	h.DB.Model(&p).Updates(input)
	res.Json(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.DB.Delete(&models.Product{}, id).Error; err != nil {
		res.Json(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка удаления"})
		return
	}
	res.Json(w, http.StatusOK, map[string]string{"message": "Товар успешно удален"})
}
