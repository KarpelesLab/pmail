package pmail

import (
	"io"
	"os/exec"
)

type SendmailSender string

var Sendmail Sender = SendmailSender("/usr/sbin/sendmail")

// Send invokes sendmail to send the specified message. from and to will be ignored
func (s SendmailSender) Send(from string, to []string, msg io.WriterTo) error {
	cmd := exec.Command(string(s))
	writer, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = msg.WriteTo(writer)
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return err
	}
	return cmd.Wait()
}
