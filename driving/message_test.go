package driving

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHTML_Successful(t *testing.T) {
	testMetadata := []map[string]string{
		{"Date": "testdate", "StartTime": "start_time"},
		{"Date": "testdate", "StartTime": "start_time"},
	}

	ts := ToHTML(testMetadata)
	assert.Equal(t, "<table border=\"1\"><tr><th>Time</th><th>Date</th></tr><tr><td>start_time</td><td>testdate</td></tr><tr><td>start_time</td><td>testdate</td></tr></table>", ts)
}

func TestToHTML_Empty(t *testing.T) {
	testMetadata := []map[string]string{}

	ts := ToHTML(testMetadata)
	assert.Equal(t, "<table border=\"1\"><tr><th>Time</th><th>Date</th></tr></table>", ts)
}
