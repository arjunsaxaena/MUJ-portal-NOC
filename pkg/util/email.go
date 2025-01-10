// pkg/util/email.go
package util

import (
	"net/smtp"
)

func SendEmail(from, to, subject, body string) error {
	password := "bjkwwhugjefvcdoa"

	// Gmail SMTP server details
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Construct the email message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	return err
}
