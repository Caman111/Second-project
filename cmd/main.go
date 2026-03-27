package main

import (
	"3-validation-api/config"
	"3-validation-api/internal/bizness"
	"3-validation-api/internal/product"
	"3-validation-api/internal/verify"
	"3-validation-api/middleware"
	"3-validation-api/pkg/db"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg := config.LoadConfig()

	conn, err := db.InitDB(cfg.DSN)
	if err != nil {
		logrus.Fatalf("Не удалось подключить к БД: %v", err)
	}
	logrus.Info("База готова, миграции прошли успешно!")

	prodHandler := &product.ProductHandler{DB: conn}
	bizHandler := &bizness.BiznessHandler{DB: conn}

	service := verify.NewService(cfg.Email, cfg.Password, cfg.Address)
	verifyHandler := verify.NewHandler(service)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", prodHandler.Create)
	mux.HandleFunc("GET /products/{id}", prodHandler.Get)
	mux.HandleFunc("PUT /products/{id}", prodHandler.Update)
	mux.HandleFunc("DELETE /products/{id}", prodHandler.Delete)

	mux.HandleFunc("POST /bizness", bizHandler.Create)
	mux.HandleFunc("GET /bizness/{id}", bizHandler.Get)
	mux.HandleFunc("PUT /bizness/{id}", bizHandler.Update)
	mux.HandleFunc("DELETE /bizness/{id}", bizHandler.Delete)

	mux.HandleFunc("/send", verifyHandler.SendEmailHandler)
	mux.HandleFunc("/verify/", verifyHandler.VerifyHandler)

	finalHandler := middleware.JSONLog(mux)

	logrus.Info("Сервер запущен на :8081")
	logrus.Fatal(http.ListenAndServe(":8081", finalHandler))
}
