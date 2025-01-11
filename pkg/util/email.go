// pkg/util/email.go
package util

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
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
	// Open the attachment file
	attachment, err := os.Open(attachmentPath)
	if err != nil {
		return fmt.Errorf("failed to open attachment: %v", err)
	}
	defer attachment.Close()

	// Read the attachment content using io.ReadAll (not ioutil.ReadAll)
	attachmentBytes, err := io.ReadAll(attachment)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %v", err)
	}

	// Prepare the email headers and body
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Add text body
	part, err := writer.CreatePart(map[string][]string{"Content-Type": {"text/plain"}})
	if err != nil {
		return err
	}
	part.Write([]byte(body))

	// Add the attachment
	attachmentPart, err := writer.CreateFormFile("attachment", attachmentPath)
	if err != nil {
		return err
	}
	attachmentPart.Write(attachmentBytes)

	// Close the multipart writer
	writer.Close()

	// Send the email
	auth := smtp.PlainAuth("", from, "bjkwwhugjefvcdoa", "smtp.gmail.com")
	return smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, buffer.Bytes())
}

// This is what the email sent by SendEmailWithAttachment looks like # debugging left

//--867df283b8c734b5e6ebc326eaefca182685b7491c615e243a3de7f57bd6 Content-Type: text/plain Dear Anushka Soni,
// Your placement application has been approved. Please find attached the No Objection Certificate (NOC) for your
//  reference. Best regards, --867df283b8c734b5e6ebc326eaefca182685b7491c615e243a3de7f57bd6 Content-Disposition:
//  form-data; name="attachment"; filename="NOC_229301124.pdf" Content-Type: application/octet-stream %PDF-1.3 3 0
//   obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 4 0 R>> endobj 4 0 obj <</Filter /FlateDecode /Length
//    370>> stream pJS]J ,v9;5 /#dB5 "(2C endstream endobj 5 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents
//	 6 0 R>> endobj 6 0 obj <</Filter /FlateDecode /Length 292>> stream H),L7| n"8L 9m2\ ;cZv) n'B" >BIHV 0$<I+ /#TCn
//	  pD4d \`u\ endstream endobj 7 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 8 0 R>> endobj 8 0 obj
//	   <</Filter /FlateDecode /Length 309>> stream h{;0  0(;)pZ` RI6, XK~i yo`T6 endstream endobj 9 0 obj <</Type /Page
//	    /Parent 1 0 R /Resources 2 0 R /Contents 10 0 R>> endobj 10 0 obj <</Filter /FlateDecode /Length 289>> stream i"$"
//		 tJ > endstream endobj 11 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 12 0 R>> endobj 12 0 obj
//		 <</Filter /FlateDecode /Length 282>> stream RseE $aMa 6:,G P:|2 n_;U lOq1p endstream endobj 13 0 obj <</Type
//		 /Page /Parent 1 0 R /Resources 2 0 R /Contents 14 0 R>> endobj 14 0 obj <</Filter /FlateDecode /Length 301>>
//		  stream H'N3 P#RK~ kX&z %Cf( endstream endobj 15 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents
//			16 0 R>> endobj 16 0 obj <</Filter /FlateDecode /Length 300>> stream zj\N o($zL endstream endobj 17 0 obj
//				<</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 18 0 R>> endobj 18 0 obj <</Filter /FlateDecode
//				 /Length 300>> stream 1O;1 n $"x S?V0 zjp2 Co`= endstream endobj 19 0 obj <</Type /Page /Parent 1 0 R
//				  /Resources 2 0 R /Contents 20 0 R>> endobj 20 0 obj <</Filter /FlateDecode /Length 284>> stream sYLx
//				  [j% a>;S endstream endobj 21 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 22 0 R>>
//				   endobj 22 0 obj <</Filter /FlateDecode /Length 290>> stream yn n]\!3 otx: endstream endobj 23 0 obj
//				   <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 24 0 R>> endobj 24 0 obj <</Filter /FlateDecode
//				    /Length 294>> stream QlKArRQ] U%]u yl+I IHmd endstream endobj 25 0 obj <</Type /Page /Parent 1 0 R /R
//					esources 2 0 R /Contents 26 0 R>> endobj 26 0 obj <</Filter /FlateDecode /Length 291>> stream 1N31 ~
//					O1eR tn?1 W}_O endstream endobj 27 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 28 0 R
//					>> endobj 28 0 obj <</Filter /FlateDecode /Length 290>> stream qwADH $rRQ *;z( VUoMu lpmz 50VU {T87 :(
//					IY endstream endobj 29 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Contents 30 0 R>> endobj 30
//					 0 obj <</Filter /FlateDecode /Length 292>> stream Gk/J$' !9qn endstream endobj 31 0 obj <</Type /Page
//					  /Parent 1 0 R /Resources 2 0 R /Contents 32 0 R>> endobj 32 0 obj <</Filter /FlateDecode /Length 300
//					  >> stream 1O31 w4[2 a_'B" endstream endobj 33 0 obj <</Type /Page /Parent 1 0 R /Resources 2 0 R /Co
//					  ntents 34 0 R>> endobj 34 0 obj <</Filter /FlateDecode /Length 256>> stream P!!n jua~ iDm/d; endstre
//					  am endobj 1 0 obj <</Type /Pages /Kids [3 0 R 5 0 R 7 0 R 9 0 R 11 0 R 13 0 R 15 0 R 17 0 R 19 0 R 2
//					  1 0 R 23 0 R 25 0 R 27 0 R 29 0 R 31 0 R 33 0 R ] /Count 16 /MediaBox [0 0 595.28 841.89] endobj 35
//					  0 obj <</Type /Font /BaseFont /Helvetica-Bold /Subtype /Type1 /Encoding /WinAnsiEncoding endobj 36 0
//					   obj <</Type /Font /BaseFont /Helvetica /Subtype /Type1 /Encoding /WinAnsiEncoding endobj 2 0 obj /P
//					   rocSet [/PDF /Text /ImageB /ImageC /ImageI] /Font << /Ff5d2de5f3a71699ae4b2d83179e62d09e6fc4126 35
//					   0 R /F0a76705d18e0494dd24cb573e53aa0a8c710ec99 36 0 R /XObject << /ColorSpace << endobj 37 0 obj /P
//					   roducer ( /CreationDate (D:20250110233010) /ModDate (D:20250110233010) endobj 38 0 obj /Type /Catal
//					   og /Pages 1 0 R /Names << /EmbeddedFiles << /Names [ ] >> endobj xref 0 39 0000000000 65535 f  0000
//					   007155 00000 n  0000007544 00000 n  0000000009 00000 n  0000000087 00000 n  0000000527 00000 n  000
//					   0000605 00000 n  0000000967 00000 n  0000001045 00000 n  0000001424 00000 n  0000001503 00000 n  00
//					   00001863 00000 n  0000001943 00000 n  0000002296 00000 n  0000002376 00000 n  0000002748 00000 n  0
//					   000002828 00000 n  0000003199 00000 n  0000003279 00000 n  0000003650 00000 n  0000003730 00000 n
//					   0000004085 00000 n  0000004165 00000 n  0000004526 00000 n  0000004606 00000 n  0000004971 00000 n
//					    0000005051 00000 n  0000005413 00000 n  0000005493 00000 n  0000005854 00000 n  0000005934 00000 n
//						  0000006297 00000 n  0000006377 00000 n  0000006748 00000 n  0000006828 00000 n  0000007345 00000
//						   n  0000007447 00000 n  0000007756 00000 n  0000007870 00000 n  trailer /Size 39 /Root 38 0 R /I
//						   nfo 37 0 R startxref 7968 %%EOF --867df283b8c734b5e6ebc326eaefca182685b7491c615e243a3de7f57bd6-
