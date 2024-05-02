package utils

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpAuthServer  = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(title string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
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

func (sender *GmailSender) SendEmail(title string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	return nil
}
