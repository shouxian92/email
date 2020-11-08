package driving

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/sendgrid/sendgrid-go"
	"github.com/shouxian92/email/service"
	"github.com/stretchr/testify/assert"
)

func TestSend_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, "all ok")
		resp.Header = http.Header{}
		return resp, nil
	})

	ctx := service.RequestContext{
		To: "me@example.com",
		Metadata: []map[string]string{
			{"Date": "testdate", "StartTime": "start_time"},
			{"Date": "testdate", "StartTime": "start_time"},
		},
	}

	c := &SendgridService{service: sendgrid.NewSendClient("")}
	err := c.Send(ctx)
	assert.Nil(t, err)
}

func TestSend_TooManyRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusTooManyRequests, "too many requests")
		resp.Header = http.Header{}
		return resp, nil
	})

	ctx := service.RequestContext{
		To: "me@example.com",
		Metadata: []map[string]string{
			{"Date": "testdate", "StartTime": "start_time"},
		},
	}

	c := &SendgridService{service: sendgrid.NewSendClient("")}
	err := c.Send(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "too many requests")
}

func TestSend_GenericError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, errors.New("some error")
	})

	ctx := service.RequestContext{
		To:       "me@example.com",
		Metadata: []map[string]string{},
	}

	c := &SendgridService{service: sendgrid.NewSendClient("")}
	err := c.Send(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), `Post "https://api.sendgrid.com/v3/mail/send": some error`)
}

func TestGetEmailInstance(t *testing.T) {
	c := GetEmailInstance()
	assert.NotNil(t, c)
	assert.NotNil(t, c.service)
}

func TestGetEmailInstance_SingleInstance(t *testing.T) {
	c1 := GetEmailInstance()
	c2 := GetEmailInstance()
	assert.NotNil(t, c1)
	assert.NotNil(t, c1.service)
	assert.Equal(t, c1, c2)
	assert.Equal(t, c1.service, c2.service)
}
