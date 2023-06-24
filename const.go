package sendmail

const (
	Mixed       = "multipart/mixed"       // when adding attachments
	Alternative = "multipart/alternative" // alternative versions, try to put text/plain first
	Related     = "multipart/related"     // when using cid: or Content-ID

	TypeHTML  = "text/html"
	TypeText  = "text/plain"
	TypeEmail = "message/rfc822"
)
