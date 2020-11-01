package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Client is an implementation of an email provider (Gmail)
type Client struct {
	service *gmail.Service
}

// Send sends an email to the intended recipient and subject with the given message
func (c *Client) Send(r string, s string, m string) (bool, error) {
	var message gmail.Message
	emailTo := "To: " + r + "\r\n"
	subject := "Subject: " + m + "\r\n"
	msg := []byte(emailTo + subject + m)
	message.Raw = base64.URLEncoding.EncodeToString(msg)
	_, err := c.service.Users.Messages.Send(emailToImpersonate, &message).Do()

	if err != nil {
		log.Fatalf("unable to send email: %v", err.(*googleapi.Error))
		return false, err
	}

	return true, nil
}

func readCredentials(ctx context.Context) (oauth2.TokenSource, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("error reading credentials json: %v", err)
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(b, gmail.GmailSendScope)

	if err != nil {
		log.Fatalf("error parsing configuration: %v", err)
		return nil, err
	}
	config.Subject = emailToImpersonate
	return config.TokenSource(ctx), nil
}

// NewInstance returns a new instance of the GmailClient
func NewInstance() (*Client, error) {
	ctx := context.Background()
	ts, err := readCredentials(ctx)
	if err != nil {
		return nil, err
	}

	srv, err := gmail.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		log.Fatalf("unable to retrieve Gmail client: %v", err)
		return nil, err
	}

	return &Client{service: srv}, nil
}
