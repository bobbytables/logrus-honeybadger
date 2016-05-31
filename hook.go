package logrushoneybadger

import (
	"github.com/Sirupsen/logrus"
	"github.com/honeybadger-io/honeybadger-go"
)

// Hook handles sending errors to honeybadger via logrus hooks
type Hook struct {
	Client *honeybadger.Client
}

// IgnoredKeys will make sure to exclude data keys in logrus entries
// in the context portion of honeybadger notifications.
type IgnoredKeys map[string]struct{}

// Add keys to ignore to the set
func (ik IgnoredKeys) Add(keys ...string) {
	for _, k := range keys {
		ik[k] = struct{}{}
	}
}

// DefaultIgnoredKeys are sane defaults for keys to ignore when sending context
// to honeybadger
var DefaultIgnoredKeys = IgnoredKeys{}

func init() {
	DefaultIgnoredKeys.Add("error")
}

// Fire implements logrus.Hook
func (h *Hook) Fire(e *logrus.Entry) error {
	var msg string

	if err, ok := e.Data["error"].(error); ok {
		msg = err.Error()
	} else {
		msg = e.Message
	}

	hbCtx := make(honeybadger.Context)
	for key, value := range e.Data {
		if _, ok := DefaultIgnoredKeys[key]; ok {
			continue
		}

		hbCtx[key] = value
	}

	_, err := h.Client.Notify(msg, hbCtx)
	return err
}

// Levels implements logrus.Hook
func (h *Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}
