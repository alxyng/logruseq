# Logruseq

A [Seq](https://getseq.net/) hook for [Logrus](https://github.com/Sirupsen/logrus)

[![GoDoc](https://godoc.org/github.com/nullseed/logruseq?status.svg)](https://godoc.org/github.com/nullseed/logruseq)

## Install

```
go get -u github.com/nullseed/logruseq
```

## Usage

```go
package main

func main() {
    log.AddHook(logruseq.NewSeqHook("http://localhost:5341"))

    log.WithFields(log.Fields{
        "animal": "walrus",
    }).Info("A walrus appears")
}
```

API keys can be specified optionally with:

```go
    logrus.AddHook(logruseq.NewSeqHook("http://localhost:5341",
                logruseq.OptionAPIKey("N1ncujiT5pYGD6m4CF0")))
```
