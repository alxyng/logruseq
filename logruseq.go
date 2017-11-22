package logruseq

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type SeqHook struct {
	host string
}

// NewSeqHook creates a Seq hook for logrus and sends log events to the host
// specified.
func NewSeqHook(host string) *SeqHook {
	return &SeqHook{
		host: host,
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

	endpoint := fmt.Sprintf("%v/api/events/raw", hook.host)
	resp, err := http.Post(endpoint, "application/vnd.serilog.clef", bytes.NewReader(data))
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
