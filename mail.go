package pmail

import (
	"bytes"
	"errors"
	"io"
	"net/mail"
	"os"
	"strings"
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

	// build default email structure: email(mixed(alternative())) (easily allows adding attachments)
	mixed := NewPart(Mixed)
	m.Body.Append(mixed)
	mixed.Append(NewPart(Alternative))

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

func (m *Mail) SetFrom(address string, name ...string) {
	if len(name) == 0 {
		m.From = &mail.Address{Address: address}
	} else {
		m.From = &mail.Address{Address: address, Name: strings.Join(name, " ")}
	}
}

func (m *Mail) AddTo(address string, name ...string) {
	if len(name) == 0 {
		m.To = append(m.To, &mail.Address{Address: address})
	} else {
		m.To = append(m.To, &mail.Address{Address: address, Name: strings.Join(name, " ")})
	}
}

func (m *Mail) AddCc(address string, name ...string) {
	if len(name) == 0 {
		m.Cc = append(m.To, &mail.Address{Address: address})
	} else {
		m.Cc = append(m.To, &mail.Address{Address: address, Name: strings.Join(name, " ")})
	}
}

func (m *Mail) AddBcc(address string, name ...string) {
	if len(name) == 0 {
		m.Bcc = append(m.To, &mail.Address{Address: address})
	} else {
		m.Bcc = append(m.To, &mail.Address{Address: address, Name: strings.Join(name, " ")})
	}
}

func (m *Mail) SetSubject(subject string) {
	m.Body.Headers.Set("Subject", subject)
}

func (m *Mail) SetBodyText(txt string) error {
	return m.SetBodyHelper([]byte(txt), "text/plain")
}

func (m *Mail) SetBodyHtml(txt string) error {
	return m.SetBodyHelper([]byte(txt), "text/html")
}

func (m *Mail) SetBodyHelper(data []byte, typ string) error {
	// set email's text body
	p := m.Body.FindType(Alternative, true)
	if p == nil {
		return errors.New("cannot use body helper without an alternative content email")
	}

	c := p.FindType(typ, false)
	if c != nil {
		// already have a part of this type, just replace the body
		c.Data = bytes.NewReader(data)
		return nil
	}

	// create part
	c = NewPart(typ)
	c.Data = bytes.NewReader(data)
	p.Append(c)

	return nil
}

// SetDate allows changing the date stored in the mail enveloppe. This should
// not be done normally, but can be useful for unit tests.
func (m *Mail) SetDate(t time.Time) {
	m.Body.Headers.Set("Date", t.Format(time.RFC1123))
}

func (m *Mail) WriteTo(w io.Writer) (int64, error) {
	m.SetTargetHeaders()

	return m.Body.WriteTo(w)
}

// SetTargetHeaders sets the various headers needed for sending the mail based on the values present in Mail
// This is called automatically when the email is sent and typically doesn't need to be manually called
func (m *Mail) SetTargetHeaders() {
	if m.From != nil {
		m.Body.Headers.SetAddressList("From", []*mail.Address{m.From})
	}
	if len(m.ReplyTo) > 0 {
		m.Body.Headers.SetAddressList("Reply-To", m.ReplyTo)
	} else {
		m.Body.Headers.Del("Reply-To")
	}
	if len(m.To) > 0 {
		m.Body.Headers.SetAddressList("To", m.To)
	} else {
		m.Body.Headers.Del("To")
	}
	if len(m.Cc) > 0 {
		m.Body.Headers.SetAddressList("Cc", m.Cc)
	} else {
		m.Body.Headers.Del("Cc")
	}
	if len(m.Bcc) > 0 {
		m.Body.Headers.SetAddressList("Bcc", m.Bcc) // we set Bcc because some email sending methods (such as sendmail) will require it
	} else {
		m.Body.Headers.Del("Bcc")
	}

	if m.MessageId == "" {
		// generate messageId (from from?)
		host := "localhost"
		if m.From != nil {
			pos := strings.LastIndexByte(m.From.Address, '@')
			if pos != -1 {
				host = m.From.Address[pos+1:]
			}
		}
		if host == "localhost" {
			if h, err := os.Hostname(); err == nil {
				host = h
			}
		}
		m.MessageId = randomBoundary() + "@" + host
	}
	m.Body.Headers.Set("Message-Id", "<"+m.MessageId+">")
}
