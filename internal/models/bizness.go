package models 

type BiznessCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	DSN      string `json:"dsn"`
	ID       uint   `json:"id"`
}
