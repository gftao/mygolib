package myLogger

import (
	"runtime"
	"path/filepath"
	"fmt"
	"github.com/sirupsen/logrus"
	"bytes"
	"strconv"
)

const (
	GLOBAL = "global" //全局name空间
)

func getEntry(name string) *logrus.Entry {
	return logger.WithFields(locate(logrus.Fields{"N": name}))
}
func locate(fields logrus.Fields) logrus.Fields {
	_, path, line, ok := runtime.Caller(3)
	if ok {
		_, file := filepath.Split(path)
		fields["F"] = fmt.Sprintf("%s:%d", file, line)
	}
	fields["G"] = fmt.Sprintf("%d", GetGID())
	return fields
}

func Debug(msg ...interface{}) {
	getEntry(GLOBAL).Debug(msg...)
}
func Debugln(msg ...interface{}) {
	getEntry(GLOBAL).Debugln(msg...)
}
func Debugf(format string, msg ...interface{}) {
	getEntry(GLOBAL).Debugf(format, msg...)
}

func Info(msg ...interface{}) {
	getEntry(GLOBAL).Info(msg...)
}
func Infoln(msg ...interface{}) {
	getEntry(GLOBAL).Infoln(msg...)
}
func Infof(format string, msg ...interface{}) {
	getEntry(GLOBAL).Infof(format, msg...)
}

func Warn(msg ...interface{}) {
	getEntry(GLOBAL).Warn(msg...)
}
func Warnln(msg ...interface{}) {
	getEntry(GLOBAL).Warnln(msg...)
}
func Warnf(format string, msg ...interface{}) {
	getEntry(GLOBAL).Warnf(format, msg...)
}
func Error(msg ...interface{}) {
	getEntry(GLOBAL).Error(msg...)
}
func Errorln(msg ...interface{}) {
	logger.Errorln(msg...)
}
func Errorf(format string, msg ...interface{}) {
	getEntry(GLOBAL).Errorf(format, msg...)
}
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}