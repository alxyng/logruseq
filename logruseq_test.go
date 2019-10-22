package logruseq

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestHookFireSetsCorrectContentTypeHeaderValue(t *testing.T) {
	hook := NewSeqHook("http://localhost:5341")
	transport := &fakeTransportRequestSaver{}
	http.DefaultClient.Transport = transport

	hook.Fire(&logrus.Entry{})

	header := "Content-Type"
	expectedContentType := "application/vnd.serilog.clef"
	actualContentType := transport.req.Header.Get(header)
	if actualContentType != expectedContentType {
		t.Errorf("incorrect value for %v, got %v, want %v",
			header, actualContentType, expectedContentType)
	}
}

func TestHookWithAPIKeyFireSetsCorrectXSeqApiKeyHeaderValue(t *testing.T) {
	expectedAPIKey := "foo"
	hook := NewSeqHook("http://localhost:5341", OptionAPIKey(expectedAPIKey))
	transport := &fakeTransportRequestSaver{}
	http.DefaultClient.Transport = transport

	hook.Fire(&logrus.Entry{})

	header := "X-Seq-ApiKey"
	actualAPIKey := transport.req.Header.Get(header)
	if actualAPIKey != expectedAPIKey {
		t.Errorf("incorrect value for %v, got %v, want %v",
			header, actualAPIKey, expectedAPIKey)
	}
}

func TestHookFireReturnsAnErrorOnRequestFailure(t *testing.T) {
	hook := NewSeqHook("http://localhost:5341")
	transport := &fakeTransportRequestError{}
	http.DefaultClient.Transport = transport

	err := hook.Fire(&logrus.Entry{})

	if err == nil {
		t.Errorf("err not nil, got %v", err)
	}
}

func TestHookFireReturnsAnErrorWhenSeqDoesNotRespondWithStatusCreated(t *testing.T) {
	body := "An error occurred"
	hook := NewSeqHook("http://localhost:5341")
	transport := &fakeTransportInternalServerError{
		body: body,
	}
	http.DefaultClient.Transport = transport
	expectedError := fmt.Errorf("error creating seq event: %v", body)

	actualError := hook.Fire(&logrus.Entry{})

	if actualError.Error() != expectedError.Error() {
		t.Errorf("incorrect error, got %v, want %v",
			actualError.Error(), expectedError.Error())
	}
}

func TestAllLevelsAreLoggedByDefault(t *testing.T) {
	levels := []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}

	hook := NewSeqHook("http://localhost:5341")

	hookLevels := hook.Levels()
	if !reflect.DeepEqual(levels, hookLevels) {
		t.Errorf("incorrect levels, got %v, want %v", hookLevels, levels)
	}
}

func TestLevelsAreCorrectlySetByOptionLevels(t *testing.T) {
	levels := []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}

	hook := NewSeqHook("http://localhost:5341", OptionLevels(levels))

	hookLevels := hook.Levels()
	if !reflect.DeepEqual(levels, hookLevels) {
		t.Errorf("incorrect levels, got %v, want %v", hookLevels, levels)
	}
}

type fakeTransportRequestSaver struct {
	req *http.Request
}

func (t *fakeTransportRequestSaver) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req

	return &http.Response{
		StatusCode: http.StatusCreated,
		Body:       bodyReadCloser{bytes.NewBuffer([]byte{})},
	}, nil
}

type fakeTransportRequestError struct{}

func (t *fakeTransportRequestError) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("example error")
}

type fakeTransportInternalServerError struct {
	body string
}

func (t *fakeTransportInternalServerError) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       bodyReadCloser{bytes.NewBuffer([]byte(t.body))},
	}, nil
}

type bodyReadCloser struct {
	io.Reader
}

func (rc bodyReadCloser) Close() error {
	return nil
}
