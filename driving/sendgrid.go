package driving

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/shouxian92/email/service"
)

var (
	// try to be singleton
	svc *SendgridService
)

// SendgridService returns a structure with an underlying sendgrid client
type SendgridService struct {
	service *sendgrid.Client
}

// Send an email to the address with the context of driving details
func (c *SendgridService) Send(ctx service.RequestContext) error {
	name := os.Getenv("EMAIL_FROM")
	email := os.Getenv("EMAIL_FROM_ADDRESS")
	fromUser := mail.NewEmail(name, email)

	to := mail.NewEmail("", ctx.To)

	msg := ToHTML(ctx.Metadata)
	message := mail.NewSingleEmail(fromUser, ctx.Subject, to, "", msg)
	resp, err := c.service.Send(message)

	if err != nil {
		log.Printf("failed to send email: %v", err)
		return err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return errors.New("too many requests")
	}

	if resp.StatusCode > http.StatusNoContent {
		log.Printf("error from sendgrid: %v", resp.Body)
		return errors.New("error encountered (" + strconv.Itoa(resp.StatusCode) + ")")
	}

	log.Printf("email sent successfully to: %v", ctx.To)
	return nil
}

// GetEmailInstance returns a new instance of a sendgrid service
func GetEmailInstance() *SendgridService {
	if svc != nil {
		return svc
	}

	svc = &SendgridService{service: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))}
	return svc
}
