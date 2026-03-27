package bizness

import (
	models "3-validation-api/internal/models"
	"3-validation-api/pkg/res"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type BiznessHandler struct {
	DB *gorm.DB
}

func (h *BiznessHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.BiznessCreate

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

func (h *BiznessHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input models.BiznessCreate

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res.Json(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	var biz models.BiznessCreate
	if err := h.DB.First(&biz, id).Error; err != nil {
		res.Json(w, http.StatusNotFound, map[string]string{"error": "Объект не найден"})
		return
	}
	if err := h.DB.Model(&biz).Updates(input).Error; err != nil {
		res.Json(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления"})
		return
	}
	res.Json(w, http.StatusOK, biz)
}

func (h *BiznessHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.DB.Delete(&models.BiznessCreate{}, id).Error; err != nil {
		res.Json(w, http.StatusInternalServerError, map[string]string{"error": "Не удалось удалить"})
		return
	}

	res.Json(w, http.StatusOK, map[string]string{"message": "Успешно удалено"})
}

func (h *BiznessHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var biz models.BiznessCreate

	if err := h.DB.First(&biz, id).Error; err != nil {
		res.Json(w, http.StatusNotFound, map[string]string{"error": "Бизнес не найден"})
		return
	}

	res.Json(w, http.StatusOK, biz)
}
