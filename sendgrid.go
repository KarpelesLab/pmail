package pmail

import (
	"encoding/base64"
	"errors"
	"io"
	"net/mail"
	"strings"

	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// AsSGMailV3 returns a sendgrid SGMailV3 object for this email
func (m *Mail) AsSGMailV3() *sgmail.SGMailV3 {
	pers := &sgmail.Personalization{
		To:   makeSGEmails(m.To),
		From: makeSGEmail(m.From),
		CC:   makeSGEmails(m.Cc),
		BCC:  makeSGEmails(m.Bcc),
	}

	// From ReplyTo[] To[] Cc[] Bcc[] Body MessageId
	res := &sgmail.SGMailV3{
		From:             makeSGEmail(m.From),
		Subject:          m.Body.Headers.Get("Subject"),
		Personalizations: []*sgmail.Personalization{pers},
	}

	// res.Content && res.Attachments
	scanSGPart(res, m.Body)

	return res
}

func scanSGPart(res *sgmail.SGMailV3, part *Part) error {
	if strings.HasPrefix(part.Type, "text/") {
		data, err := io.ReadAll(part.Data)
		if err != nil {
			return err
		}
		res.AddContent(sgmail.NewContent(part.Type, string(data)))
		return nil
	} else if part.IsContainer() {
		for _, sub := range part.Children {
			scanSGPart(res, sub)
		}
		return nil
	} else if part.Data != nil {
		// attachment
		data, err := io.ReadAll(part.Data)
		if err != nil {
			return err
		}
		attach := &sgmail.Attachment{
			Content: base64.StdEncoding.EncodeToString(data),
			Type:    part.Type,
			// TODO Filename, etc
		}
		res.AddAttachment(attach)
		return nil
	}
	return errors.New("unsupported mime type")
}

func makeSGEmail(addr *mail.Address) *sgmail.Email {
	return &sgmail.Email{
		Name:    addr.Name,
		Address: addr.Address,
	}
}

func makeSGEmails(addrs []*mail.Address) []*sgmail.Email {
	res := make([]*sgmail.Email, len(addrs))
	for n, a := range addrs {
		res[n] = makeSGEmail(a)
	}
	return res
}

func convertSGHeaders(h Header) map[string]string {
	res := make(map[string]string)
	for k, v := range h {
		switch k {
		case "Subject", "Content-Type":
			// do nothing
		default:
			res[k] = v[0]
		}
	}
	return res
}
