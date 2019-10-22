package main

import (
	"github.com/nullseed/logruseq"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.AddHook(logruseq.NewSeqHook("http://localhost:5341"))

	// Optionally, the hook can be used with an API key:
	// log.AddHook(logruseq.NewSeqHook("http://localhost:5341",
	// 	logruseq.OptionAPIKey("N1ncujiT5pYGD6m4CF0")))

	// Optionally, which levels to log can be specified:
	// log.AddHook(logruseq.NewSeqHook("http://localhost:5341",
	// 	logruseq.OptionLevels([]logrus.Level{
	// 		logrus.WarnLevel,
	// 		logrus.ErrorLevel,
	// 		logrus.FatalLevel,
	// 		logrus.PanicLevel,
	// 	})))

	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
