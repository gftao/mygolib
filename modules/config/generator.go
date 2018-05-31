package config

import (
	"errors"
	"fmt"
	"os"
)

const CONFIGPREFIX = "config="

var (
	initFlg bool = false
)

func InitModule() error {

	var filename string = ""

	if len(os.Args) < 2 {
		fmt.Println("需要配置参数文件：config=***")
		return errors.New("配置文件未配置")
	}

	for i := 1; i < len(os.Args); i++ {
		if CONFIGPREFIX == os.Args[i][:7] {
			filename = os.Args[i][7:]
			break
		}
	}

	if filename == "" {
		fmt.Println("需要配置参数文件：config=***")
		return errors.New("配置文件未配置")
	}

	var err error
	instance, err = LoadConfig(filename)
	if err != nil {
		return err
	}

	fmt.Println("配置文件加载成功:", filename)
	initFlg = true

	return nil
}

func InitModuleByParams(filename string) error {

	var err error
	instance, err = LoadConfig(filename)
	if err != nil {
		return err
	}

	fmt.Println("配置文件加载成功:", filename)
	initFlg = true

	return nil
}

func HasConfigInit() bool {
	return initFlg
}
