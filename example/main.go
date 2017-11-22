package main

import (
	"github.com/nullseed/logruseq"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.AddHook(logruseq.NewSeqHook("http://localhost:5341"))

	// Or optionally use the hook with an API key:
	// log.AddHook(logruseq.NewSeqHook("http://localhost:5341",
	// 	logruseq.OptionAPIKey("N1ncujiT5pYGD6m4CF0")))

	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
