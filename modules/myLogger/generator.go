package myLogger

import (
	"github.com/gogap/logrus_mate"
	"os"
	"errors"
	_ "github.com/gogap/logrus_mate/hooks/file"
	"fmt"
	"mygolib/modules/config"
	"github.com/sirupsen/logrus"
)

var initFlg bool
var logger *logrus.Logger

func InitLoggers() error {
	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}
	config.SetSection("main")
	runMode := config.StringDefault("runmode", "")
	logFileName := config.StringDefault("logFileName", "")

	if logFileName == "" || runMode == "" {
		return errors.New("日志配置文件未配置")
	}
	os.Setenv("RUN_MODE", runMode)
	mate, err := logrus_mate.NewLogrusMate(logrus_mate.ConfigFile(logFileName))
	if err != nil {
		return err
	}

	newLoger := logrus.New()
	if err = mate.Hijack(newLoger, "guft"); err != nil {
		fmt.Println(err)
		return err
	}
	logger = newLoger
	initFlg = true

	return nil
}

func HasLoggerInit() bool {
	return initFlg
}

func GetNameLog(name string) (l *logrus.Logger) {
	if HasLoggerInit() {
		l = logger.WithFields(logrus.Fields{"N": name}).Logger
	} else {
		l = logrus.New().WithFields(logrus.Fields{"N": name}).Logger
	}
	return l
}
