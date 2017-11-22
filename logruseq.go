package logruseq

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// SeqHook sends logs to Seq via HTTP.
type SeqHook struct {
	endpoint string
	apiKey   string
}

// NewSeqHook creates a Seq hook for logrus which can send log events to the
// host specified, for example:
//     logruseq.NewSeqHook("http://localhost:5341")
// Optionally, the hook can be used with an API key, for example:
//     logruseq.NewSeqHook("http://localhost:5341", logruseq.OptionAPIKey("N1ncujiT5pYGD6m4CF0"))
func NewSeqHook(host string, opts ...func(*SeqHookOptions)) *SeqHook {
	sho := &SeqHookOptions{}

	for _, opt := range opts {
		opt(sho)
	}

	endpoint := fmt.Sprintf("%v/api/events/raw", host)

	return &SeqHook{
		endpoint: endpoint,
		apiKey:   sho.apiKey,
	}
}

// Fire sends a log entry to Seq.
func (hook *SeqHook) Fire(entry *logrus.Entry) error {
	formatter := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "@mt",
			logrus.FieldKeyLevel: "@l",
			logrus.FieldKeyTime:  "@t",
		},
	}

	data, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", hook.endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/vnd.serilog.clef")

	if hook.apiKey != "" {
		req.Header.Add("X-Seq-ApiKey", hook.apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		data, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}

		return fmt.Errorf("error creating seq event: %v", string(data))
	}

	return nil
}

// Levels returns the levels for which Fire will be called for. These are Debug,
// Info, Warn, Error and Fatal. Verbose level events (featured in Seq) and
// Panic level events (featured in Logrus) are not handled.
func (hook *SeqHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}
}

// SeqHookOptions collects non-default Seq hook options.
type SeqHookOptions struct {
	apiKey string
}

// OptionAPIKey sets the Seq API key option.
func OptionAPIKey(apiKey string) func(opts *SeqHookOptions) {
	return func(opts *SeqHookOptions) {
		opts.apiKey = apiKey
	}
}
