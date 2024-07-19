package pmail

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"strings"

	"github.com/KarpelesLab/rndpass"
)

type Part struct {
	Type     string  // text/plain, etc
	Children []*Part // only if multipart/* type
	Data     io.ReadCloser
	GetBody  func() (io.ReadCloser, error)
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
		p.Boundary = rndpass.Code(24, rndpass.RangeFull)
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
	return strings.HasPrefix(p.Type, "multipart/")
}

func (p *Part) IsEmail() bool {
	return p.Type == TypeEmail
}

func (p *Part) IsContainer() bool {
	return p.IsMultipart() || p.Type == TypeEmail
}

func (p *Part) IsEmpty() bool {
	if len(p.Children) == 1 {
		return p.Children[0].IsEmpty()
	}
	return len(p.Children) == 0 && p.Data == nil
}

func (p *Part) Append(c *Part) {
	p.Children = append(p.Children, c)
}

// WriteTo writes the part to the given output
func (p *Part) WriteTo(w io.Writer) (int64, error) {
	wc := &writeCounter{W: w}
	w = wc

	hdrs := p.Headers
	isEmail := p.IsEmail()

	for len(p.Children) == 1 {
		// if only 1 child, move down and merge headers
		p = p.Children[0]
		hdrs = hdrs.Merge(p.Headers)
	}

	if len(p.Children) > 0 {
		// enforce content type & boundary
		hdrs.Set("Content-Type", p.Type+"; boundary="+p.Boundary)
	}

	// Write headers
	w.Write(hdrs.Encode())
	w.Write([]byte{'\r', '\n'})

	isMultipart := strings.HasPrefix(hdrs.Get("Content-Type"), "multipart/")

	// If this is a multipart email, write the "this is a mime message" line
	if isEmail && isMultipart {
		w.Write([]byte("This is a message in Mime Format.  If you see this, your mail reader does not support this format.\r\n\r\n"))
	}

	if enc := p.convertEncoder(w); enc != nil {
		defer enc.Close()
		w = enc
	}
	if len(p.Children) == 0 {
		// simply copy
		if p.Data == nil {
			if p.GetBody != nil {
				var err error
				p.Data, err = p.GetBody()
				if err != nil {
					return wc.C, err
				}
			} else {
				return wc.C, ErrPartHasNoBody
			}
		}
		// take p.Data, make it nil instead so we don't use it twice (and GetBody gets called next time, if there is a next time)
		fp := p.Data
		p.Data = nil
		defer fp.Close()
		_, err := io.Copy(w, fp)
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

func (p *Part) readBody() ([]byte, error) {
	if p.Data == nil {
		if p.GetBody != nil {
			var err error
			p.Data, err = p.GetBody()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, ErrPartHasNoBody
		}
	}
	// take p.Data, make it nil instead so we don't use it twice (and GetBody gets called next time, if there is a next time)
	fp := p.Data
	p.Data = nil
	defer fp.Close()
	return io.ReadAll(fp)
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

func (p *Part) FindType(typ string, recurse bool) *Part {
	// search direct children
	for _, c := range p.Children {
		if c.Type == typ {
			return c
		}
	}
	if !recurse {
		return nil
	}
	// search containers children
	for _, c := range p.Children {
		if c.IsContainer() {
			if f := c.FindType(typ, true); f != nil {
				return f
			}
		}
	}
	return nil
}
