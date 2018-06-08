package run

import (
	"fmt"
	"mygolib/gerror"
	"path/filepath"
	"runtime"
	"mygolib/modules/myLogger"
	"github.com/sirupsen/logrus"
)

type BaseWorker struct {
	Id       uint32
	NodeName string
	OrderId  string
}

func (t *BaseWorker) Init() gerror.IError {

	return nil
}

func (t *BaseWorker) SetOrderId(orderId string) {
	t.OrderId = orderId
}

func (t *BaseWorker) SetSysOrderId(sysOrderId string) {
	t.OrderId = sysOrderId
}

func (t *BaseWorker) getEntry() *logrus.Entry {
	l := myLogger.GetNameLog(t.NodeName)

	return l.WithFields(locate(logrus.Fields{"I": t.Id, "O": t.OrderId}))
}

func locate(fields logrus.Fields) logrus.Fields {
	_, path, line, ok := runtime.Caller(3)
	if ok {
		_, file := filepath.Split(path)
		fields["F"] = fmt.Sprintf("%s:%d", file, line)
	}
	fields["G"] = fmt.Sprintf("%d", myLogger.GetGID())
	return fields
}

func (t BaseWorker) Info(msg ...interface{}) {
	t.getEntry().Infoln(msg...)
}

func (t BaseWorker) Debug(msg ...interface{}) {
	t.getEntry().Debugln(msg...)
}

func (t BaseWorker) Error(msg ...interface{}) {
	t.getEntry().Errorln(msg...)
}

func (t BaseWorker) Warn(msg ...interface{}) {
	t.getEntry().Warnln(msg...)
}

func (t BaseWorker) Infof(format string, msg ...interface{}) {
	t.getEntry().Infof(format, msg...)
}

func (t BaseWorker) Debugf(format string, msg ...interface{}) {
	t.getEntry().Debugf(format, msg...)
}

func (t BaseWorker) Errorf(format string, msg ...interface{}) {
	t.getEntry().Errorf(format, msg...)
}

func (t BaseWorker) Warnf(format string, msg ...interface{}) {
	t.getEntry().Warnf(format, msg...)
}
