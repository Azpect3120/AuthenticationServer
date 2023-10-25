package model

import (
	"net/smtp"
)

func SendEmail (to, subject, content string) (*Email, *Error) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := "user.authentication.server12@gmail.com"
	smtpPassword := "jjyu ddmu deag xwpe"

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		smtpUsername,
		[]string{to},
		[]byte("Subject: " + subject + "\r\n" + content),
	)

	if err != nil {
		return nil, &Error{ Status: 500, Message: err.Error() }
	} else {
		return &Email{ To: to, Subject: subject, Content: content }, nil
	}

}
