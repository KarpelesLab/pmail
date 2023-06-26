package pmail_test

import (
	"bytes"
	"net/mail"
	"testing"
	"time"

	"github.com/KarpelesLab/pmail"
)

func TestTextMail(t *testing.T) {
	// generate an email
	m := pmail.New()
	m.From = &mail.Address{Address: "test@example.com", Name: "Test"}
	m.SetDate(time.Unix(1687756384, 0).UTC())
	m.MessageId = "test1@localhost"
	m.AddTo("bob@example.com", "Bob Test")
	m.SetSubject("Hello Bob")
	m.SetBodyText("Hello Bob,\n\nCan you look at this?")

	buf := &bytes.Buffer{}
	m.WriteTo(buf)

	expect := []byte(`Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain
Date: Mon, 26 Jun 2023 05:13:04 UTC
From: "Test" <test@example.com>
Message-Id: <test1@localhost>
Mime-Version: 1.0
Subject: Hello Bob
To: "Bob Test" <bob@example.com>

Hello Bob,

Can you look at this?`)
	expect = bytes.ReplaceAll(expect, []byte{'\n'}, []byte{'\r', '\n'})

	if !bytes.Equal(expect, buf.Bytes()) {
		t.Errorf("mail not as expected.\nexpected:\n%s\noutput:\n%s", expect, buf.Bytes())
	}

	//log.Printf("Email:\n%s", buf.Bytes())
	//log.Printf("hash = %x", sha256.Sum256(buf.Bytes()))

}

func TestAltEmail(t *testing.T) {
	m := pmail.New()
	m.SetFrom("test@example.com", "Test")
	m.SetDate(time.Unix(1687756384, 0).UTC())
	m.MessageId = "test2@localhost"
	m.AddTo("bob@example.com", "Bob Test")
	m.SetSubject("Hello Bob")
	m.SetBodyText("Hello Bob,\r\n\r\nCan you look at this?")
	m.SetBodyHtml("<p>Hello Bob,</p>\r\n<p>Can you look at this?</p>")

	m.Body.FindType(pmail.Alternative, true).Boundary = "test123456"

	buf := &bytes.Buffer{}
	m.WriteTo(buf)

	expect := []byte(`Content-Transfer-Encoding: 7bit
Content-Type: multipart/alternative; boundary=test123456
Date: Mon, 26 Jun 2023 05:13:04 UTC
From: "Test" <test@example.com>
Message-Id: <test2@localhost>
Mime-Version: 1.0
Subject: Hello Bob
To: "Bob Test" <bob@example.com>

This is a message in Mime Format.  If you see this, your mail reader does not support this format.


--test123456
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain

Hello Bob,

Can you look at this?
--test123456
Content-Transfer-Encoding: quoted-printable
Content-Type: text/html

<p>Hello Bob,</p>
<p>Can you look at this?</p>
--test123456--
`)
	expect = bytes.ReplaceAll(expect, []byte{'\n'}, []byte{'\r', '\n'})

	if !bytes.Equal(expect, buf.Bytes()) {
		t.Errorf("mail not as expected.\nexpected:\n%s\noutput:\n%s", expect, buf.Bytes())
	}

	//log.Printf("Email:\n%s", buf.Bytes())
}
