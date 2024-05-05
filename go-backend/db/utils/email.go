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
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, emailAddress, emailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  emailAddress,
		fromEmailPassword: emailPassword,
	}
}

func (sender *GmailSender) SendEmail(title string, content string, mimeType string, to []string, cc []string, bcc []string, attachFiles []string, attachmentData [][]byte) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.HTML = []byte(content)

	for i, f := range attachmentData {
		attachReader := bytes.NewReader(f)
		_, err := e.Attach(attachReader, attachFiles[i], mimeType)
		if err != nil {
			return err
		}
	}

	return nil
}
