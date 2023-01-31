package database

import (
	"github.com/sirupsen/logrus"
	"xorm.io/xorm/log"
)

type dbLogger struct {
	level   log.LogLevel
	showSQL bool
}

func (ll *dbLogger) Debug(v ...interface{}) {
	logrus.Debug(v...)
}

func (ll *dbLogger) Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func (ll *dbLogger) Error(v ...interface{}) {
	logrus.Error(v...)
}

func (ll *dbLogger) Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func (ll *dbLogger) Info(v ...interface{}) {
	logrus.Info(v...)
}

func (ll *dbLogger) Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func (ll *dbLogger) Warn(v ...interface{}) {
	logrus.Warn(v...)
}

func (ll *dbLogger) Warnf(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
}

func (ll *dbLogger) Level() log.LogLevel {
	return ll.level
}

func (ll *dbLogger) SetLevel(l log.LogLevel) {
	ll.level = l
}

func (ll *dbLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		ll.showSQL = true
		return
	}
	ll.showSQL = show[0]
}

func (ll *dbLogger) IsShowSQL() bool {
	return ll.showSQL
}
