package verify

import (
	"fmt"
	"net"
	"net/smtp"
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

func (s *Service) SaveHash(hash, userEmail string) {
	store[hash] = userEmail
}

func (s *Service) VerifyHash(hash string) (string, bool) {
	email, ok := store[hash]
	if ok {
		delete(store, hash)
	}
	return email, ok
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
