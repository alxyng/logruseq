# Logruseq

A Seq hook for [Logrus](https://github.com/Sirupsen/logrus) :walrus:

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
