package logruseq

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type SeqHook struct {
	host string
}

func NewSeqHook(host string) *SeqHook {
	return &SeqHook{
		host: host,
	}
}

func (hook *SeqHook) Fire(entry *logrus.Entry) error {
	formatter := logrus.JSONFormatter{
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

func (hook *SeqHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}
}
