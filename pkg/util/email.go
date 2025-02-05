package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
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

func downloadFile(url, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch file: received status %d", resp.StatusCode)
	}

	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func SendEmailWithAttachment(from, to, subject, body, appPassword, nocPath string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, appPassword, smtpHost)

	nocURL := fmt.Sprintf("http://localhost:8002/files/NOC/%s", filepath.Base(nocPath))
	tempFilePath := fmt.Sprintf("./temp_%s", filepath.Base(nocPath))

	err := downloadFile(nocURL, tempFilePath)
	if err != nil {
		log.Printf("Error downloading NOC: %v", err)
		return fmt.Errorf("failed to download NOC: %v", err)
	}
	defer os.Remove(tempFilePath)

	attachment, err := os.ReadFile(tempFilePath)
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

	attachmentPart, _ := writer.CreatePart(map[string][]string{
		"Content-Disposition":       {"attachment; filename=\"" + filepath.Base(nocPath) + "\""},
		"Content-Type":              {"application/pdf; name=\"" + filepath.Base(nocPath) + "\""},
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
