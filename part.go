package sendmail

import (
	"io"
	"net/textproto"
	"strings"
)

type Part struct {
	Type     string  // text/plain, etc
	Children []*Part // only if multipart/* type
	Data     io.Reader
	Headers  textproto.MIMEHeader
	Boundary string
	Encoding byte
}

func NewPart(typ string) *Part {
	p := &Part{
		Type:    typ,
		Headers: make(textproto.MIMEHeader),
	}
	if p.IsMultipart() {
		p.Boundary = randomBoundary()
		p.Encoding = 0
		p.Headers.Set("Content-Transfer-Encoding", "7bit")
	} else if strings.HasPrefix(p.Type, "text/") {
		// use quoted printable
		p.Headers.Set("Content-Transfer-Encoding", "quoted-printable")
		p.Encoding = 'q'
	} else {
		p.Headers.Set("Content-Transfer-Encoding", "base64")
		p.Encoding = 'b'
	}
	return p
}

func (p *Part) IsMultipart() bool {
	return strings.HasPrefix(p.Type, "multipart/")
}

func (p *Part) IsContainer() bool {
	return p.IsMultipart() || p.Type == TypeEmail
}
