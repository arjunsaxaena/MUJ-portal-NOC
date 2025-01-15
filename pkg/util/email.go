package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/smtp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func SendEmail(from, to, subject, body string) error {
	// Ensure this is your Gmail app-specific password
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

func downloadFileFromS3(s3URL string, bucket string) ([]byte, error) {
	key := s3URL[len("https://muj-student-data.s3.amazonaws.com/"):]

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}
	s3Client := s3.New(sess)

	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}
	defer resp.Body.Close()

	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 file content: %v", err)
	}

	return fileContent, nil
}

func SendEmailWithAttachment(from string, to string, subject string, body string, s3URL string) error {
	// Log the input parameters
	fmt.Printf("Sending email from: %s to: %s with subject: %s and attachment: %s\n", from, to, subject, s3URL)

	attachment, err := downloadFileFromS3(s3URL, "muj-student-data")
	if err != nil {
		fmt.Printf("Error fetching attachment: %v\n", err)
		return fmt.Errorf("failed to fetch attachment from S3: %v", err)
	}

	// Base64 encode the attachment
	encodedAttachment := base64.StdEncoding.EncodeToString(attachment)

	fmt.Println("Attachment successfully encoded.")

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
	message += "Content-Disposition: attachment; filename=\"attachment.pdf\"\n"
	message += "Content-Transfer-Encoding: base64\n\n"
	message += encodedAttachment + "\n"
	message += "--boundary--"

	fmt.Println("Email message successfully created.")

	auth := smtp.PlainAuth("", from, "bjkwwhugjefvcdoa", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(message))
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully.")
	return nil
}
