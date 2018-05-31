package myLogger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/gogap/logrus_mate"
	"github.com/gogap/config"
)

const TIMEFORMAT = "2006-01-02T15:04:05.999999999"

var logField = []string{
	"N", //结点名
	"G", //GOROUTINE id
	"F", //文件名，行数
	"M", //说明信息
}

func init() {
	logrus_mate.RegisterFormatter("myFormatter", NewMyLogFormatter)
}

type MyFormatterConfig struct {
	InHook  bool
	Address string `json:"address"`
}

func (f *MyFormatterConfig) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	if !f.InHook {
		data["time"] = entry.Time.Format(TIMEFORMAT)
		data["level"] = entry.Level.String()
	}
	data["M"] = entry.Message

	return f.formatLog(data)
}
func (f *MyFormatterConfig) formatLog(data logrus.Fields) ([]byte, error) {
	//按如下格式输出日志
	var bf bytes.Buffer
	if !f.InHook {
		bf.WriteByte('[')
		fmt.Fprint(&bf, data["time"])
		bf.WriteByte(']')
		bf.WriteByte(' ')
		bf.WriteByte('[')
		fmt.Fprint(&bf, data["level"])
		bf.WriteByte(']')
		bf.WriteByte(' ')
	}
	for _, f := range logField {
		if v, ok := data[f]; ok {
			bf.WriteString(f)
			bf.WriteByte('[')
			fmt.Fprint(&bf, v)
			bf.WriteByte(']')
			bf.WriteByte(' ')
		}
	}

	return bf.Bytes(), nil
}

func NewMyLogFormatter(config config.Configuration) (formatter logrus.Formatter, err error) {
	formatter = &MyFormatterConfig{InHook: false}
	return
}

func GetFields(entry *logrus.Entry, level logrus.Level) (logrus.Fields, bool, error) {
	data := make(logrus.Fields, len(entry.Data)+3)

	if entry.Level > level {
		return nil, false, nil
	}

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	data["time"] = entry.Time.Format(TIMEFORMAT)
	data["level"] = entry.Level.String()
	data["M"] = entry.Message

	return data, true, nil
}
