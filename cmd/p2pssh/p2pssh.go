package main

import (
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/p2pssh/cmd"
)

func main() {

	// initial log formatting; this setting is updated after the daemon configuration is loaded.
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
		FullTimestamp:   true,
	})

	if err := cmd.Execute(); err != nil {
		// logrus.Println(err)
	}
	logrus.Debug("pssh exit")
}
