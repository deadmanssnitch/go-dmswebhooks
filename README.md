# Receiver for Dead Man's Snitch Webhooks

[![GoDoc](https://godoc.org/github.com/deadmanssnitch/go-dmswebhooks?status.svg)](http://godoc.org/github.com/deadmanssnitch/go-dmswebhooks)

## Basic Usage:
```go
// NewHandler takes a callback function which is called when alerts are received
handler := dmswebhooks.NewHandler(
  func(alert *dmswebhooks.Alert) error {
    if(alert.Type == dmswebhooks.TypeSnitchReporting) {
      fmt.Println("🎉")
    }

    return nil
  }
)

```

## [Hipchat Example](https://github.com/deadmanssnitch/go-hipchat-example)

For a more complete (and fully functional example) we've built a basic Hipchat
alerter.
