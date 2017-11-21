package logruseq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type resource struct {
	Events []event
}

type event struct {
	Timestamp       string
	Level           string
	MessageTemplate string
	Properties      logrus.Fields
}

type SeqHook struct {
	host string
}

func NewSeqHook(host string) *SeqHook {
	return &SeqHook{
		host: host,
	}
}

func (hook *SeqHook) Fire(entry *logrus.Entry) error {
	r := resource{
		Events: []event{
			event{
				Timestamp:       entry.Time.UTC().Format(time.RFC3339Nano),
				Level:           entry.Level.String(),
				MessageTemplate: entry.Message,
				Properties:      entry.Data,
			},
		},
	}
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%v/api/events/raw", hook.host)
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(data))
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
