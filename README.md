# sendgridrus - SendGrid hook into Logrus
sendgridrus is a package exposing a SendGrid hook into [Logrus](https://github.com/Sirupsen/logrus).
If your machine doesn't have access to an SMTP port or you just want to use SendGrid, this is
a viable option for logging.

By default, sendgridrus hooks will fire on all entries with levels Warning, Error, or Panic.

## Installation

sendgridrus can be installed using `go get`:
```go
go get github.com/jnschaeffer/sendgridrus
```

## Usage

The simplest way to use sendgridrus is to add it to a Logrus Logger.

```go
import (
	"github.com/Sirupsen/logrus"
	"github.com/jnschaeffer/sendgridrus"
)

func main() {
	log := logrus.New()
	hook := sendgridrus.NewHook("SENDGRID API KEY", "SERVICE NAME", "FROM ADDRESS", "TO ADDRESS")
	log.Hooks.Add(hook)
}
```
