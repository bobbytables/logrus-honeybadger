package logrushoneybadger

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/stretchr/testify/assert"
)

type testBackend struct {
	sync.Mutex
	features []honeybadger.Feature
	payloads []honeybadger.Payload
}

func newTestBackend() *testBackend {
	return &testBackend{
		features: make([]honeybadger.Feature, 0),
		payloads: make([]honeybadger.Payload, 0),
	}
}

func (tb *testBackend) Notify(f honeybadger.Feature, p honeybadger.Payload) error {
	tb.Lock()
	tb.features = append(tb.features, f)
	tb.payloads = append(tb.payloads, p)
	tb.Unlock()

	return nil
}

func TestHookDispatchesToHoneybadger(t *testing.T) {
	b, c, l := setup()

	l.Error("this is an error")
	c.Flush()

	assert.Len(t, b.payloads, 1)
	notice := convertToNotice(b.payloads[0])
	assert.Equal(t, "this is an error", notice.ErrorMessage)
}

func TestHookWithErrorDispatchesToHoneybadger(t *testing.T) {
	b, c, l := setup()

	l.WithError(fmt.Errorf("yo crap dun broke")).Error("this is an error")
	c.Flush()

	assert.Len(t, b.payloads, 1)
	notice := convertToNotice(b.payloads[0])
	assert.Equal(t, "yo crap dun broke", notice.ErrorMessage)
}

func TestHookWithFieldsIncludesContextToHoneybadger(t *testing.T) {
	b, c, l := setup()

	l.WithFields(logrus.Fields{"host": "paperwalls", "user": "bobbytables"}).Error("this is an error")
	c.Flush()

	assert.Len(t, b.payloads, 1)
	notice := convertToNotice(b.payloads[0])
	assert.Equal(t, "this is an error", notice.ErrorMessage)
	assert.Equal(t, "paperwalls", notice.Context["host"])
	assert.Equal(t, "bobbytables", notice.Context["user"])
	assert.Empty(t, notice.Context["error"])
}

func convertToNotice(p honeybadger.Payload) *honeybadger.Notice {
	return p.(*honeybadger.Notice)
}

func setup() (*testBackend, *honeybadger.Client, *logrus.Logger) {
	b := newTestBackend()
	c := honeybadger.New(honeybadger.Configuration{Backend: b})
	hook := &Hook{Client: c}
	l := logrus.New()
	l.Out = ioutil.Discard
	l.Hooks.Add(hook)

	return b, c, l
}
