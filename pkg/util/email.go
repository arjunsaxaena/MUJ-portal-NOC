// pkg/util/email.go
package util

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(from, to, subject, body string) error {
	password := "bjkwwhugjefvcdoa"

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

func SendEmailWithAttachment(from string, to string, subject string, body string, attachmentPath string) error {
	// Read the attachment file
	attachment, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %v", err)
	}

	// Base64 encode the attachment
	encodedAttachment := base64.StdEncoding.EncodeToString(attachment)

	// Create the email body
	message := ""
	message += fmt.Sprintf("From: %s\n", from)
	message += fmt.Sprintf("To: %s\n", to)
	message += fmt.Sprintf("Subject: %s\n", subject)
	message += "MIME-Version: 1.0\n"
	message += "Content-Type: multipart/mixed; boundary=boundary\n\n"

	// Add the plain text body
	message += "--boundary\n"
	message += "Content-Type: text/plain; charset=utf-8\n\n"
	message += body + "\n\n"

	// Add the attachment
	message += "--boundary\n"
	message += "Content-Type: application/pdf\n"
	message += fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", attachmentPath)
	message += "Content-Transfer-Encoding: base64\n\n"
	message += encodedAttachment + "\n"
	message += "--boundary--"

	auth := smtp.PlainAuth("", from, "bjkwwhugjefvcdoa", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
