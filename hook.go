package sendgridrus

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/sendgrid/sendgrid-go"
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
	sg          *sendgrid.SGClient
	serviceName string
	fromAddr    string
	toAddr      string
	levels      []logrus.Level
}

// NewHook creates a new Logrus hook into SendGrid.
func NewHook(sgAPIKey, serviceName, fromAddr, toAddr string) *Hook {
	sg := sendgrid.NewSendGridClientWithApiKey(sgAPIKey)
	levels := []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel, logrus.PanicLevel}

	return &Hook{
		sg:          sg,
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
// Each mesasge body contains the message text, log entry time, and fields.
func (h *Hook) Fire(e *logrus.Entry) error {
	var bodyBuf bytes.Buffer

	errParse := tmpl.Execute(&bodyBuf, e)
	if errParse != nil {
		return errParse
	}

	subject := fmt.Sprintf("%s[%s]: %s", h.serviceName, e.Level, e.Message)

	message := sendgrid.NewMail()
	message.AddTo(h.toAddr)
	message.SetSubject(subject)
	message.SetText(string(bodyBuf.Bytes()))
	message.SetFrom(h.fromAddr)

	errSend := h.sg.Send(message)
	if errSend != nil {
		return errSend
	}

	return nil
}

func foo() {
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		fmt.Println("Environment variable SENDGRID_API_KEY is undefined. Did you forget to source sendgrid.env?")
		os.Exit(1)
	}
	sg := sendgrid.NewSendGridClientWithApiKey(sendgridKey)
	message := sendgrid.NewMail()
	message.AddTo("community@sendgrid.com")
	message.AddToName("SendGrid Community Dev Team")
	message.SetSubject("SendGrid Testing")
	message.SetText("WIN")
	message.SetFrom("you@yourdomain.com")
	if r := sg.Send(message); r == nil {
		fmt.Println("Email sent!")
	} else {
		fmt.Println(r)
	}
}
