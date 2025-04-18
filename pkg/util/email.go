package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"path/filepath"
)

func SendEmail(from, to, subject, body, appPassword string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, appPassword, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	log.Printf("Attempting to send email...")
	log.Printf("From: %s, To: %s", from, to)
	log.Printf("Subject: %s", subject)
	log.Printf("Body: %s", body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("email sent successfully to: %s", to)
	return nil
}

func SendEmailWithAttachment(from, to, subject, body, appPassword, nocFileName string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, appPassword, smtpHost)

	// Resolve the absolute path of the attachment
	// nocPath, err := filepath.Abs(filepath.Join("..", "uploads", "NOC", nocFileName))
	// if err != nil {
	// 	log.Printf("Error resolving file path: %v", err)
	// 	return fmt.Errorf("failed to resolve file path: %v", err)
	// }

	nocPath := filepath.Join("/app/uploads", "NOC", nocFileName)

	// Check if the file exists
	if _, err := os.Stat(nocPath); os.IsNotExist(err) {
		log.Printf("Attachment file does not exist: %s", nocPath)
		return fmt.Errorf("attachment file does not exist: %s", nocPath)
	}

	attachment, err := os.ReadFile(nocPath)
	if err != nil {
		log.Printf("Error reading attachment: %v", err)
		return fmt.Errorf("failed to read attachment: %v", err)
	}

	var emailBuffer bytes.Buffer
	writer := multipart.NewWriter(&emailBuffer)

	boundary := writer.Boundary()
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n",
		from, to, subject, boundary,
	)
	emailBuffer.WriteString(headers)

	bodyPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=utf-8"},
	})
	bodyWriter := quotedprintable.NewWriter(bodyPart)
	bodyWriter.Write([]byte(body))
	bodyWriter.Close()

	filename := filepath.Base(nocFileName)

	attachmentPart, _ := writer.CreatePart(map[string][]string{
		"Content-Disposition":       {"attachment; filename=\"" + filename + "\""},
		"Content-Type":              {"application/pdf; name=\"" + filename + "\""},
		"Content-Transfer-Encoding": {"base64"},
	})

	encodedAttachment := base64.StdEncoding.EncodeToString(attachment)
	for i := 0; i < len(encodedAttachment); i += 76 {
		end := i + 76
		if end > len(encodedAttachment) {
			end = len(encodedAttachment)
		}
		attachmentPart.Write([]byte(encodedAttachment[i:end] + "\r\n"))
	}

	writer.Close()

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, emailBuffer.Bytes())
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email with attachment sent successfully to: %s", to)
	return nil
}
