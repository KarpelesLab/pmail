package sendmail

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Part struct {
	Type     string  // text/plain, etc
	Children []*Part // only if multipart/* type
	Data     io.Reader
	Headers  Header
	Boundary string
	Encoding byte
}

func NewPart(typ string) *Part {
	p := &Part{
		Type:    typ,
		Headers: make(Header),
	}
	if p.IsMultipart() {
		p.Boundary = randomBoundary()
		p.Encoding = 0
		p.Headers.Set("Content-Transfer-Encoding", "7bit")
		p.Headers.Set("Content-Type", p.Type+"; boundary="+p.Boundary)
	} else if strings.HasPrefix(p.Type, "text/") {
		// use quoted printable
		p.Headers.Set("Content-Transfer-Encoding", "quoted-printable")
		p.Headers.Set("Content-Type", p.Type)
		p.Encoding = 'q'
	} else {
		p.Headers.Set("Content-Transfer-Encoding", "base64")
		p.Headers.Set("Content-Type", p.Type)
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

func (p *Part) GetHeaders(exclude ...string) []byte {
	// build an exclude map
	excl := make(map[string]bool)
	for _, v := range exclude {
		excl[strings.ToLower(v)] = true
	}

	buf := &bytes.Buffer{}

	for k, v := range p.Headers {
		_, skip := excl[strings.ToLower(k)]
		if skip {
			continue
		}

		switch k {
		case "Subject", "From", "To", "Cc":
			for _, s := range v {
				smartEncodeHeader(buf, k, s)
			}
		default:
			for _, s := range v {
				fmt.Fprintf(buf, "%s: %s\r\n", k, s)
			}
		}
	}
	return buf.Bytes()
}

func smartEncodeHeader(buf *bytes.Buffer, k string, v string) {
	// TODO
	fmt.Fprintf(buf, "%s: %s\r\n", k, v)
}
