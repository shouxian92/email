package service

// Client represents a particular email service
type Client interface {
	// Send sends an email using the underlying service, it takes in a recipient, subject and message as inputs
	Send(ctx RequestContext) error
}

// RequestContext is used for sending the email message in context of the email service
type RequestContext struct {
	Type     string
	To       string
	Subject  string
	Message  string
	Metadata []map[string]string
}
