package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shouxian92/email/service"
	"github.com/stretchr/testify/assert"
)

type testClient struct{}

func (c *testClient) Send(r service.RequestContext) error {
	return nil
}

type fourTwoNineClient struct{}

func (c *fourTwoNineClient) Send(r service.RequestContext) error {
	return errors.New("too many requests")
}

type internalServerClient struct{}

func (c *internalServerClient) Send(r service.RequestContext) error {
	return errors.New("unknown error")
}

func TestSendHandler_Driving_Success(t *testing.T) {
	clients := map[string]service.Client{
		"driving": &testClient{},
	}

	ts := httptest.NewServer(http.HandlerFunc(SendHandler(clients)))
	defer ts.Close()

	payload := "{\"type\":\"driving\", \"to\":\"test@example.com\", \"metadata\": [{\"StartTime\": \"sometime\", \"Date\": \"somedate\"}]}"

	res, err := ts.Client().Post(ts.URL, "application/json", strings.NewReader(payload))
	if err != nil {
		assert.Fail(t, "failed to POST to SendHandler: %v", err)
	}

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestSendHandler_BadRequests(t *testing.T) {
	testCases := []struct {
		name          string
		payload       string
		expectedError string
	}{
		{"NonExistentEmailClient", "{\"type\":\"foo\", \"to\":\"test@example.com\", \"metadata\": [{\"StartTime\": \"sometime\", \"Date\": \"somedate\"}]}", "{\"error\":\"unrecognized email request type\"}"},
		{"UnableToParseRequestBody", "not json", "{\"error\":\"failed to parse request body\"}"},
	}

	clients := map[string]service.Client{
		"driving": &testClient{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(SendHandler(clients)))
			defer ts.Close()

			payload := tc.payload

			res, err := ts.Client().Post(ts.URL, "application/json", strings.NewReader(payload))
			if err != nil {
				assert.Fail(t, "failed to POST to SendHandler: %v", err)
			}

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			b, _ := ioutil.ReadAll(res.Body)
			assert.Equal(t, tc.expectedError, strings.Trim(string(b), "\n"))
		})
	}
}

func TestSendHandler_ErrorStatusCodes(t *testing.T) {
	testCases := []struct {
		name               string
		client             service.Client
		expectedStatusCode int
	}{
		{"TooManyRequests", &fourTwoNineClient{}, http.StatusTooManyRequests},
		{"InternalServerError", &internalServerClient{}, http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientName := "foo"
			clients := map[string]service.Client{
				clientName: tc.client,
			}

			ts := httptest.NewServer(http.HandlerFunc(SendHandler(clients)))
			defer ts.Close()

			payload := "{\"type\":\"" + clientName + "\", \"to\":\"test@example.com\", \"metadata\": [{\"StartTime\": \"sometime\", \"Date\": \"somedate\"}]}"

			res, err := ts.Client().Post(ts.URL, "application/json", strings.NewReader(payload))
			if err != nil {
				assert.Fail(t, "failed to POST to SendHandler: %v", err)
			}

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestDecodeRequestBody_Success(t *testing.T) {
	payload := []byte("{\"type\":\"foo\", \"to\":\"test@example.com\", \"message\":\"kekw\"}")
	ctx, err := decodeRequestBody(payload)

	assert.Nil(t, err)
	assert.Equal(t, "foo", ctx.Type)
	assert.Equal(t, "test@example.com", ctx.To)
	assert.Equal(t, "kekw", ctx.Message)
}

func TestDecodeRequestBody_Errors(t *testing.T) {
	testCases := []struct {
		name          string
		payload       string
		expectedError string
	}{
		{"MissingTo", "{\"type\":\"foo\", \"message\":\"kekw\"}", "'to' is required"},
		{"MissingType", "{\"to\":\"hannibal\", \"message\":\"kekw\"}", "'type' is required"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload := []byte(tc.payload)
			ctx, err := decodeRequestBody(payload)
			assert.NotEmpty(t, ctx)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}
