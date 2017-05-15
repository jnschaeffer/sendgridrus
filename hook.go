package sendgridrus

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var tmpl *template.Template

func init() {
	tmplText := `{{.Level.String}}: {{.Message}}

Logged at: {{.Time}}
{{with .Data}}

Fields:
{{range $key, $val := .}}{{$key}}: {{$val}}
{{end}}
{{end}}`

	var err error
	tmpl, err = template.New("sendgrid").Parse(tmplText)
	if err != nil {
		log.Fatal(err)
	}
}

// Hook represents a Logrus hook to SendGrid. Users should have a valid API
// key from SendGrid before using this hook.
type Hook struct {
	sgAPIKey    string
	serviceName string
	fromAddr    string
	toAddr      string
	levels      []logrus.Level
}

// NewHook creates a new Logrus hook into SendGrid.
func NewHook(sgAPIKey, serviceName, fromAddr, toAddr string) *Hook {
	levels := []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel, logrus.PanicLevel}

	return &Hook{
		sgAPIKey:    sgAPIKey,
		serviceName: serviceName,
		fromAddr:    fromAddr,
		toAddr:      toAddr,
		levels:      levels,
	}
}

// Levels returns the levels for this hook. By default these are Warn, Error,
// and Panic. Use (*Hook).SetLevels to change the levels for this hook.
func (h *Hook) Levels() []logrus.Level {
	return h.levels
}

// SetLevels sets the logging levels for this hook.
func (h *Hook) SetLevels(levels []logrus.Level) {
	h.levels = levels
}

// Fire sends a message on SendGrid from the hook's From user to the hook's To
// user. Messages have a subject with the format "SERVICE_NAME[LEVEL]: MESSAGE".
// Each message body contains the message text, log entry time, and fields.
func (h *Hook) Fire(e *logrus.Entry) error {
	var bodyBuf bytes.Buffer

	errParse := tmpl.Execute(&bodyBuf, e)
	if errParse != nil {
		return errParse
	}

	from := mail.NewEmail(h.serviceName, h.fromAddr)
	msg := e.Message
	if len(msg) > 50 {
		msg = msg[:47] + "..."
	}
	subject := fmt.Sprintf("%s[%s]: %s", h.serviceName, e.Level, msg)
	to := mail.NewEmail("", h.toAddr)
	content := mail.NewContent("text/plain", string(bodyBuf.Bytes()))
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(h.sgAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, errSend := sendgrid.API(request)
	if errSend != nil {
		return errSend
	}

	return nil
}
