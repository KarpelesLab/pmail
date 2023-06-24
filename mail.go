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

func (m *Mail) SetFrom(from *mail.Address) {
	m.From = from
	m.Body.Headers.Set("From", from.String())
}
