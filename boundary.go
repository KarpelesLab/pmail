package sendmail

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func randomBoundary() string {
	r := make([]byte, 28)
	_, err := io.ReadFull(rand.Reader, r)
	if err != nil {
		// ?????
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(r)
}
