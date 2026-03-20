package verify

import (
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"

	"github.com/jordan-wright/email"
)

var store = map[string]string{}

type Service struct {
	ConfigEmail string
	ConfigPass  string
	ConfigAddr  string
}

func NewService(email, pass, addr string) *Service {
	return &Service{
		ConfigEmail: email,
		ConfigPass:  pass,
		ConfigAddr:  addr,
	}
}

func (s *Service) SaveHash(hash string, email string) {
	var users = make(map[string]string)
	content, err := os.ReadFile("users.json")
	if err == nil {
		json.Unmarshal(content, &users)
	}
	users[hash] = email

	file, _ := json.MarshalIndent(users, "", "  ")
	_ = os.WriteFile("users.json", file, 0644)

	fmt.Println("Хеш успешно привязан к почте в users.json")
}

func (s *Service) VerifyHash(hash string) (string, bool) {
	file, err := os.ReadFile("users.json")
	if err != nil {
		return "", false
	}
	var users map[string]string
	if err := json.Unmarshal(file, &users); err != nil {
		return "", false
	}
	email, ok := users[hash]
	if ok {
		delete(users, hash)
		updatedFile, _ := json.MarshalIndent(users, "", "  ")
		_ = os.WriteFile("users.json", updatedFile, 0644)
		return email, true
	}

	return "", false
}

func (s *Service) SendEmail(to, link string) error {
	e := email.NewEmail()
	e.From = s.ConfigEmail
	e.To = []string{to}
	e.Subject = "Подтверждение Email"
	e.Text = []byte(fmt.Sprintf("Перейдите по ссылке для подтверждения: %s", link))
	host, _, _ := net.SplitHostPort(s.ConfigAddr)

	auth := smtp.PlainAuth("", s.ConfigEmail, s.ConfigPass, host)
	return e.Send(s.ConfigAddr, auth)
}
func (s *Service) GenerateHash() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
