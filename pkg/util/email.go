package util

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(from, to, subject, body string, appPassword string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", from, appPassword, smtpHost)

	// Construct the email message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	// Log the email details for debugging
	log.Printf("Attempting to send email...")
	log.Printf("From: %s, To: %s", from, to)
	log.Printf("Subject: %s", subject)
	log.Printf("Body: %s", body)

	// Sending email via SMTP
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		// Log the error in case of failure
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	// Log successful email sending
	log.Printf("Email sent successfully to: %s", to)
	return nil
}
