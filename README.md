# package logrushoneybadger

This package allows hooking into logrus + honeybadger. It will send errors and context to Honeybadger.

## Quickstart

```go
client := honeybadger.New(honeybadger.Configuration{})
hook := &logrushoneybadger.Hook{Client: client}

logrus.AddHook(hook)
```

Anytime you log with a level of `error`, `fatal`, or `panic`, Logrus will now
dispatch the error to Honeybadger.

### Using the default honeybadger client

```go
hook := &logrushoneybadger.Hook{Client: honeybadger.DefaultClient}
logrus.AddHook(hook)
```
