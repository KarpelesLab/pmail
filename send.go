package pmail

import "io"

type Sender interface {
	Send(from string, to []string, msg io.WriterTo) error
}

// Send sends the email using the given Sender.
func (m *Mail) Send(s Sender) error {
	if !m.IsValid() {
		return ErrInvalidEmail
	}

	to := make([]string, 0, len(m.To))
	for _, a := range m.To {
		to = append(to, a.Address)
	}

	return s.Send(m.From.Address, to, m)
}
