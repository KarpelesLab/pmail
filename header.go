package sendmail

import (
	"net/mail"
	"net/textproto"
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

func (h Header) Date() (time.Time, error) {
	// parse & return date, if any
	v := h.Get("Date")
	if v == "" {
		return time.Time{}, mail.ErrHeaderNotPresent
	}

	return mail.ParseDate(v)
}

func (h Header) AddressList(key string) ([]*mail.Address, error) {
	hdr := h.Get(key)
	if hdr == "" {
		return nil, mail.ErrHeaderNotPresent
	}
	return mail.ParseAddressList(hdr)
}
