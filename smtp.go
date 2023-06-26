package pmail

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
)

// RelaySender sends emails to a relaying server
type RelaySender struct {
	Host       string
	Port       int
	Auth       smtp.Auth
	TLSConfig  *tls.Config
	RequireTLS bool
}

func NewDialer(host string, port int, username, password string) *RelaySender {
	return &RelaySender{Host: host, Port: port, Auth: smtp.PlainAuth("", username, password, host)}
}

// Send connects to the relay server and sends the email
func (r *RelaySender) Send(from string, to []string, msg io.WriterTo) error {
	cl, err := smtp.Dial(fmt.Sprintf("%s:%d", r.Host, r.Port))
	if err != nil {
		return err
	}
	defer cl.Close()

	if stls, _ := cl.Extension("STARTTLS"); stls {
		tlscfg := &tls.Config{ServerName: r.Host}
		if r.TLSConfig != nil {
			tlscfg = r.TLSConfig
		}
		err = cl.StartTLS(tlscfg)
		if err != nil {
			return err
		}
	} else if r.RequireTLS {
		return fmt.Errorf("cannot send email: configuration requires TLS but server %s does not support it", r.Host)
	}

	if r.Auth != nil {
		err = cl.Auth(r.Auth)
		if err != nil {
			return err
		}
	}

	err = cl.Mail(from)
	if err != nil {
		return err
	}
	for _, t := range to {
		err = cl.Rcpt(t)
		if err != nil {
			return err
		}
	}

	// send data
	wr, err := cl.Data()
	if err != nil {
		return err
	}
	_, err = msg.WriteTo(wr)
	if err != nil {
		return err
	}
	err = wr.Close()
	if err != nil {
		return err
	}

	cl.Quit()
	return nil
}
