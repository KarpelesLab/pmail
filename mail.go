package sendmail

import (
	"net/mail"
	"os"
	"time"
)

type Mail struct {
	From      *mail.Address
	ReplyTo   []*mail.Address
	To        []*mail.Address
	Cc        []*mail.Address
	Bcc       []*mail.Address
	Body      *Part
	MessageId string
}

func New() *Mail {
	m := &Mail{}
	m.Body = NewPart(TypeEmail)
	m.Body.Headers.Set("MIME-Version", "1.0")

	host := "localhost"
	if h, err := os.Hostname(); err == nil {
		host = h
	}
	m.MessageId = randomBoundary() + "@" + host
	m.Body.Headers.Set("Message-Id", "<"+m.MessageId+">")
	m.Body.Headers.Set("Date", time.Now().Format(time.RFC1123)) // RFC 5322 compatible, it seems
	return m
}

// IsValid returns true if the email is valid and can be sent
func (m *Mail) IsValid() bool {
	if m.From == nil {
		return false
	}
	if len(m.To) == 0 {
		return false
	}
	if m.Body.IsEmpty() {
		return false
	}
	return true
}

// SetTargetHeaders sets the various headers needed for sending the mail based on the values present in Mail
// This is called automatically when the email is sent and typically doesn't need to be manually called
func (m *Mail) SetTargetHeaders() {
	m.Body.Headers.SetAddressList("From", []*mail.Address{m.From})
	m.Body.Headers.SetAddressList("Reply-To", m.ReplyTo)
	m.Body.Headers.SetAddressList("To", m.To)
	m.Body.Headers.SetAddressList("Cc", m.Cc)
	m.Body.Headers.SetAddressList("Bcc", m.Bcc) // we set Bcc because some email sending methods (such as sendmail) will require it
}
