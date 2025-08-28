package delivery

// EmailMessage is the message to be sent via email
type EmailMessage struct {
	From      string
	To        []string
	Subject   string
	Body      string
	IsHTML    bool
	FilePaths []string
}
