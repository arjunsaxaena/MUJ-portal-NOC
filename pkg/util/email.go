// package util

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"mime/multipart"
// 	"mime/quotedprintable"
// 	"net/smtp"
// 	"os"
// 	"path/filepath"
// )

// func SendEmail(from, to, subject, body, appPassword string) error {
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	auth := smtp.PlainAuth("", from, appPassword, smtpHost)

// 	msg := []byte("To: " + to + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body + "\r\n")

// 	log.Printf("Attempting to send email...")
// 	log.Printf("From: %s, To: %s", from, to)
// 	log.Printf("Subject: %s", subject)
// 	log.Printf("Body: %s", body)

// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
// 	if err != nil {
// 		log.Printf("Error sending email: %v", err)
// 		return fmt.Errorf("failed to send email: %v", err)
// 	}

// 	log.Printf("email sent successfully to: %s", to)
// 	return nil
// }

// // func downloadFile(url, savePath string) error {
// // 	resp, err := http.Get(url)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to fetch file: %v", err)
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusOK {
// // 		return fmt.Errorf("failed to fetch file: received status %d", resp.StatusCode)
// // 	}

// // 	outFile, err := os.Create(savePath)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to create file: %v", err)
// // 	}
// // 	defer outFile.Close()

// // 	_, err = io.Copy(outFile, resp.Body)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to save file: %v", err)
// // 	}

// // 	return nil
// // }

// func SendEmailWithAttachment(from, to, subject, body, appPassword, nocFileName string) error {
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"
// 	auth := smtp.PlainAuth("", from, appPassword, smtpHost)

// 	nocPath := filepath.Join("..", "uploads", nocFileName)

// 	attachment, err := os.ReadFile(nocPath)
// 	if err != nil {
// 		log.Printf("Error reading attachment: %v", err)
// 		return fmt.Errorf("failed to read attachment: %v", err)
// 	}

// 	var emailBuffer bytes.Buffer
// 	writer := multipart.NewWriter(&emailBuffer)

// 	boundary := writer.Boundary()
// 	headers := fmt.Sprintf(
// 		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n",
// 		from, to, subject, boundary,
// 	)
// 	emailBuffer.WriteString(headers)

// 	bodyPart, _ := writer.CreatePart(map[string][]string{
// 		"Content-Type": {"text/plain; charset=utf-8"},
// 	})
// 	bodyWriter := quotedprintable.NewWriter(bodyPart)
// 	bodyWriter.Write([]byte(body))
// 	bodyWriter.Close()

// 	filename := filepath.Base(nocFileName)

// 	attachmentPart, _ := writer.CreatePart(map[string][]string{
// 		"Content-Disposition":       {"attachment; filename=\"" + filename + "\""},
// 		"Content-Type":              {"application/pdf; name=\"" + filename + "\""},
// 		"Content-Transfer-Encoding": {"base64"},
// 	})

// 	encodedAttachment := base64.StdEncoding.EncodeToString(attachment)
// 	for i := 0; i < len(encodedAttachment); i += 76 {
// 		end := i + 76
// 		if end > len(encodedAttachment) {
// 			end = len(encodedAttachment)
// 		}
// 		attachmentPart.Write([]byte(encodedAttachment[i:end] + "\r\n"))
// 	}

// 	writer.Close()

// 	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, emailBuffer.Bytes())
// 	if err != nil {
// 		log.Printf("Error sending email: %v", err)
// 		return fmt.Errorf("failed to send email: %v", err)
// 	}

// 	log.Printf("Email with attachment sent successfully to: %s", to)
// 	return nil
// }

package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendEmail(senderEmail, recipientEmail, subject, body string) error {
	srv, err := getGmailService(senderEmail)
	if err != nil {
		return err
	}

	messageStr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", senderEmail, recipientEmail, subject, body)
	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(messageStr)),
	}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendEmailWithAttachment(senderEmail, recipientEmail, subject, body, nocFileName string) error {
	srv, err := getGmailService(senderEmail)
	if err != nil {
		return fmt.Errorf("failed to get Gmail service: %w", err)
	}

	nocPath, err := filepath.Abs(filepath.Join("..", "uploads", "NOC", nocFileName))
	if err != nil {
		return fmt.Errorf("failed to resolve file path: %w", err)
	}
	log.Println("Resolved NOC path:", nocPath)

	fileData, err := ioutil.ReadFile(nocPath)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %w", err)
	}
	filename := filepath.Base(nocPath)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	boundary := writer.Boundary()
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n",
		senderEmail, recipientEmail, subject, boundary,
	)
	buffer.WriteString(headers)

	bodyPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qp := quotedprintable.NewWriter(bodyPart)
	qp.Write([]byte(body))
	qp.Close()

	attachmentPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"application/pdf; name=" + filename},
		"Content-Disposition":       {"attachment; filename=" + filename},
		"Content-Transfer-Encoding": {"base64"},
	})

	encodedAttachment := base64.StdEncoding.EncodeToString(fileData)
	for i := 0; i < len(encodedAttachment); i += 76 {
		end := i + 76
		if end > len(encodedAttachment) {
			end = len(encodedAttachment)
		}
		attachmentPart.Write([]byte(encodedAttachment[i:end] + "\r\n"))
	}

	writer.Close()

	rawMessage := base64.URLEncoding.EncodeToString(buffer.Bytes())
	message := &gmail.Message{Raw: rawMessage}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("failed to send email with attachment: %w", err)
	}

	log.Println("Email sent successfully to:", recipientEmail)
	return nil
}

func getGmailService(senderEmail string) (*gmail.Service, error) {
	ctx := context.Background()

	b, err := os.ReadFile("../tokens/credential-gmail-api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials.json: %w", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials.json: %w", err)
	}

	tokenFile := fmt.Sprintf("../tokens/token_%s.json", strings.ReplaceAll(strings.ReplaceAll(senderEmail, "@", "_"), ".", "_"))
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file for %s: %w", senderEmail, err)
	}

	var tokenObj oauth2.Token
	err = json.Unmarshal(token, &tokenObj)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	client := config.Client(ctx, &tokenObj)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %w", err)
	}

	return srv, nil
}
