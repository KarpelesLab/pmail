package pmail

// RFC 5322 explicitly states for body: CR and LF MUST only occur together as CRLF;
// they MUST NOT appear independently in the body.

func verifycrlf(buf []byte) bool {
	pb := byte(0)

	for _, b := range buf {
		if b == '\n' && pb != '\r' {
			// \n not preceded by \r = bad
			return false
		}
		if pb == '\r' && b != '\n' {
			// \r not followed by \n = bad
			return false
		}
		pb = b
	}
	if pb == '\r' {
		// \r at end of buffer = bad
		return false
	}
	return true
}

func fixcrlf(buf []byte) []byte {
	// replace any CR or LF by CRLF, unless it is a CRLF
	out := make([]byte, 0, len(buf)+(len(buf)/50)) // approximately too much

	pb := byte(0)

	for _, b := range buf {
		switch {
		case b == '\n' && pb != '\r':
			out = append(out, '\r', '\n')
		case pb == '\r' && b != '\n':
			out = append(out, '\r', '\n', b)
		case b == '\r':
			// do nothing
		case pb == '\r' && b == '\n':
			// good linebreak, append it
			out = append(out, '\r', '\n')
		default:
			out = append(out, b)
		}
		pb = b
	}
	return out
}
