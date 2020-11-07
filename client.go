package main

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Client is an implementation of an email provider (Gmail)
type Client struct {
	service *sendgrid.Client
}

// Send sends an email to the intended recipient and subject with the given message
func (c *Client) Send(r string, s string, m string) (bool, error) {
	fromUser := mail.NewEmail(senderName, senderEmail)
	to := mail.NewEmail("", r)
	message := mail.NewSingleEmail(fromUser, s, to, "", m)
	_, err := c.service.Send(message)

	if err != nil {
		log.Printf("failed to send email: %v", err)
		return false, err
	}

	return true, nil
}

// NewInstance returns a new instance of the email client
func NewInstance() (*Client, error) {
	return &Client{service: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))}, nil
}
