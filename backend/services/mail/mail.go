package mail

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/cmd/config"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
	"github.com/xhit/go-simple-mail/v2"
)

type MailService struct {
	mailService *mail.SMTPServer
}

func NewMailService(cfg config.MailServiceConfig) *MailService {
	var mailService MailService

	server := mail.NewSMTPClient()
	server.Host = cfg.Host
	server.Port, _ = strconv.Atoi(cfg.Port)
	server.Username = cfg.Username
	server.Password = cfg.Password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	mailService.mailService = server

	return &mailService
}

func (m *MailService) templateResetPassword(username, newRandPassword string) string {
	return fmt.Sprintf("Hello, %s\n\nYour password has been reset to %s. Please log in to your account and change the password as soon as possible.", username, newRandPassword)
}

func (m *MailService) ResetPassword(user *core.User, newRandPassword string) (serializer.Serializer, error) {
	client, err := m.mailService.Connect()
	if err != nil {
		return serializer.NewSerializer(true, "unable to initialize email client to notify user", nil), err
	}

	msg := m.templateResetPassword(user.DisplayName, newRandPassword)

	email := mail.NewMSG()
	email.SetFrom(m.mailService.Username).AddTo(user.Email).SetReplyTo(m.mailService.Username).SetSubject("Password reset for The NJC Reader")
	email.SetBody(mail.TextPlain, msg)

	err = email.Send(client)
	if err != nil {
		return serializer.NewSerializer(true, "unable to send password reset email to user", nil), err
	}

	return serializer.NewSerializer(false, "successfully reset password. please check your email for the new temporary password", nil), nil
}
