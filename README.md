# Logrus Seq

A Seq hook for [Logrus](https://github.com/Sirupsen/logrus)

## Install

```
go get -u github.com/nullseed/logrus-seq
```

## Usage

```go
package main

func main() {
    log.AddHook(logrus_seq.NewSeqHook("http://localhost:5341"))

    log.WithFields(log.Fields{
        "animal": "walrus",
    }).Info("A walrus appears")
}
```
