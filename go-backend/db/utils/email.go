package utils

import (
	"bytes"
	"fmt"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpAuthServer  = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(title string, content string, mimeType string, to []string, cc []string, bcc []string, attachFiles []string, attachmentData [][]byte) error
}

type GmailSender struct {
	name                string
	from_email_address  string
	from_email_password string
}

func NewGmailSender(name, email_address, email_password string) EmailSender {
	return &GmailSender{
		name:                name,
		from_email_address:  email_address,
		from_email_password: email_password,
	}
}

func (sender *GmailSender) SendEmail(title string, content string, mimeType string, to []string, cc []string, bcc []string, attachFiles []string, attachmentData [][]byte) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.from_email_address)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.HTML = []byte(content)

	for i, f := range attachmentData {
		attachReader := bytes.NewReader(f)
		e.Attach(attachReader, attachFiles[i], mimeType)
	}

	return nil
}
