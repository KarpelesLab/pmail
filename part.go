package pmail

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
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
	if strings.HasPrefix(typ, "multipart/") {
		p.Boundary = randomBoundary()
		p.Encoding = 0
		p.Headers.Set("Content-Transfer-Encoding", "7bit")
		p.Headers.Set("Content-Type", p.Type+"; boundary="+p.Boundary)
	} else if typ == TypeEmail {
		// do nothing, wait for more info
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
	return strings.HasPrefix(p.Headers.Get("Content-Type"), "multipart/")
}

func (p *Part) IsEmail() bool {
	return p.Type == TypeEmail
}

func (p *Part) IsMultipartEmail() bool {
	return p.IsEmail() && p.IsMultipart()
}

func (p *Part) IsContainer() bool {
	return p.IsMultipart() || p.Type == TypeEmail
}

func (p *Part) IsEmpty() bool {
	return len(p.Children) == 0 && p.Data == nil
}

// WriteTo writes the part to the given output
func (p *Part) WriteTo(w io.Writer) (int64, error) {
	wc := &writeCounter{W: w}
	w = wc

	// Write headers
	w.Write(p.Headers.Encode())
	w.Write([]byte{'\r', '\n'})

	// If this is a multipart email, write the "this is a mime message" line
	if p.IsMultipartEmail() {
		w.Write([]byte("This is a message in Mime Format.  If you see this, your mail reader does not support this format.\r\n\r\n"))
	}

	if enc := p.convertEncoder(w); enc != nil {
		defer enc.Close()
		w = enc
	}
	if len(p.Children) == 0 {
		// simply copy
		_, err := io.Copy(w, p.Data)
		return wc.C, err
	}

	// for each children...
	for _, child := range p.Children {
		// boundary start
		fmt.Fprintf(w, "\r\n--%s\r\n", p.Boundary)
		_, err := child.WriteTo(w)
		if err != nil {
			return wc.C, err
		}
	}
	// boundary end
	fmt.Fprintf(w, "\r\n--%s--\r\n", p.Boundary)

	return wc.C, nil
}

func (p *Part) convertEncoder(w io.Writer) io.WriteCloser {
	switch p.Encoding {
	case 'q':
		return quotedprintable.NewWriter(w)
	case 'b':
		breaker := &lineBreaker{out: w}
		return base64.NewEncoder(base64.StdEncoding, breaker)
	default:
		return nil
	}
}
