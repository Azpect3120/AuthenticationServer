package model

import (
	"github.com/joho/godotenv"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, content string) (*Email, *Error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := os.Getenv("smtp_email")
	smtpPassword := os.Getenv("smtp_password")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpUsername,
		[]string{to},
		[]byte("Subject: "+subject+"\r\n"+content),
	)

	if err != nil {
		return nil, &Error{Status: 500, Message: err.Error()}
	} else {
		return &Email{To: to, Subject: subject, Content: content}, nil
	}

}
