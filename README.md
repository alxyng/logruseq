# Logruseq

A [Seq](https://getseq.net/) hook for [Logrus](https://github.com/Sirupsen/logrus)

[![Build Status](https://travis-ci.org/nullseed/logruseq.svg?branch=master)](https://travis-ci.org/nullseed/logruseq)
[![Go Report Card](https://goreportcard.com/badge/github.com/nullseed/logruseq)](https://goreportcard.com/report/github.com/nullseed/logruseq)
[![GoDoc](https://godoc.org/github.com/nullseed/logruseq?status.svg)](https://godoc.org/github.com/nullseed/logruseq)

## Install

```
go get -u github.com/nullseed/logruseq
```

## Usage

```go
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
```
