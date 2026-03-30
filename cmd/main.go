package main

import (
	"3-validation-api/config"
	"3-validation-api/internal/auth"
	"3-validation-api/internal/bizness"
	"3-validation-api/internal/product"
	"3-validation-api/internal/verify"
	"3-validation-api/middleware"
	"3-validation-api/pkg/db"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg := config.LoadConfig()

	conn, err := db.InitDB(cfg.DSN)
	if err != nil {
		logrus.Fatalf("Не удалось подключить к БД: %v", err)
	}

	prodHandler := &product.ProductHandler{DB: conn}
	bizHandler := &bizness.BiznessHandler{DB: conn}

	service := verify.NewService(cfg.Email, cfg.Password, cfg.Address)
	verifyHandler := verify.NewHandler(service)

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

	authRepo := auth.NewAuthRepository()
	authHandler := &auth.AuthHandler{Repo: authRepo}

	mux.HandleFunc("POST /auth/login", authHandler.Login())
	mux.HandleFunc("POST /auth/verify", authHandler.Verify())
	mux.HandleFunc("/auth/me", middleware.AuthMiddleware(authHandler.GetProfile))

	finalHandler := middleware.JSONLog(mux)

	fmt.Println("Сервер запущен на :8081")
	logrus.Fatal(http.ListenAndServe(":8081", finalHandler))
}
