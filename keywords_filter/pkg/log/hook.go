package log

import (
	"log"

	"github.com/sirupsen/logrus"
)

type errorhook struct {
}

func (*errorhook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

func (*errorhook) Fire(entry *logrus.Entry) error {
	log.Println(entry.Message, entry.Data)
	return nil
}
