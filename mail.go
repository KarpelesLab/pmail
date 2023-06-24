package sendmail

import "net/mail"

type Mail struct {
	From    *mail.Address
	ReplyTo []*mail.Address
	To      []*mail.Address
	Cc      []*mail.Address
	Bcc     []*mail.Address
	Body    *Part
}

func New() *Mail {
	m := &Mail{}
	m.Body = NewPart(TypeEmail)
	return m
}

func (m *Mail) SetFrom(from *mail.Address) {
	m.From = from
	m.Body.Headers.Set("From", from.String())
}
