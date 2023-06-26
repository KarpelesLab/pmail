package pmail

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/textproto"
	"sort"
	"time"
)

type Header map[string][]string

func (h Header) Add(key, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	h[key] = append(h[key], value)
}

func (h Header) Set(key, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	h[key] = []string{value}
}

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	v := h[textproto.CanonicalMIMEHeaderKey(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (h Header) Del(key string) {
	delete(h, textproto.CanonicalMIMEHeaderKey(key))
}

func (h Header) Date() (time.Time, error) {
	// parse & return date, if any
	v := h.Get("Date")
	if v == "" {
		return time.Time{}, mail.ErrHeaderNotPresent
	}

	return mail.ParseDate(v)
}

func (h Header) SetAddressList(key string, value []*mail.Address) {
	// reverse of AddressList
	buf := &bytes.Buffer{}

	for n, a := range value {
		if n > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(a.String())
	}
	h.Set(key, buf.String())
}

func (h Header) AddressList(key string) ([]*mail.Address, error) {
	hdr := h.Get(key)
	if hdr == "" {
		return nil, mail.ErrHeaderNotPresent
	}
	return mail.ParseAddressList(hdr)
}

func (h Header) Encode(exclude ...string) []byte {
	// build an exclude map
	excl := make(map[string]bool)
	for _, v := range exclude {
		excl[textproto.CanonicalMIMEHeaderKey(v)] = true
	}

	buf := &bytes.Buffer{}

	// sort keys to ensure we always produce the same output
	keys := make([]string, 0, len(h))
	for k := range h {
		_, skip := excl[k] // we assume k is already canonical
		if skip {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := h[k]

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

// Merge will duplicate the header object and add another object
func (h Header) Merge(h2 Header) Header {
	n := make(Header)
	for k, v := range h {
		n[k] = v
	}
	for k, v := range h2 {
		n[k] = v
	}
	return n
}
