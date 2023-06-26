[![GoDoc](https://godoc.org/github.com/KarpelesLab/pmail?status.svg)](https://godoc.org/github.com/KarpelesLab/pmail)

# pmail

go-based email sending lib, including html emails, attachements, etc and with various email sending methods

# Sample usage

```go
	m := pmail.New()
	m.SetSubject("test")
	m.SetFrom("test@localhost")
	m.AddTo("bob@localhost")
	m.SetBodyText("Hi\nThis is an email!\n")

	m.Send(pmail.Sendmail) // on linux, if sendmail is configured
```
