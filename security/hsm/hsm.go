package hsm

import (
	"golib/modules/logs"
	"net"
	"golib/modules/config"
	"strings"
	"fmt"
	"os"
)

type HsmServer struct {
	HsmName  string
	HsmType  string
	HsmStat  bool
	ConnNum  int
	HsmAddr  string
	ConnChan chan *net.TCPConn
}

const (
	HSM_SJL05 = "SJL05"
	HSM_SJL06 = "SJL06"
	HSM_CONF  = "HSM"
)

//全局加密机对象
var GlbSvr *HsmServer

func HsmInit() {
	//装载配置
	config.SetSection(HSM_CONF)
	GlbSvr = new(HsmServer)
	GlbSvr.HsmName = config.StringDefault("HsmName", "SJL05")
	GlbSvr.HsmAddr = config.StringDefault("HsmAddr", "")
	GlbSvr.ConnNum = config.IntDefault("ConnNum", 1)
	GlbSvr.HsmType = config.StringDefault("HsmType",HSM_SJL05)
	GlbSvr.ConnChan = make(chan *net.TCPConn, GlbSvr.ConnNum*2+10)

	logCfg := config.StringDefault("LogConf", "console,{}")
	LogInit(logCfg)

	switch GlbSvr.HsmType {
	case HSM_SJL05:
		Hsm05Svr = GlbSvr
		err := Hsm05Svr.Init()
		if err != nil {
			logs.Error("Hsm05Svr.Init Error:[%s] ", err)
			return
		}
		logs.Info("Hsm05Svr Init Success!")
		GlbSvr = Hsm05Svr
	default:
		logs.Error("非法加密机类型，加密机初始化失败[%s]",GlbSvr.HsmType )
		break
	}
	return
}

func HsmClose() {
	GlbSvr.Close()
}


func LogInit(cfg string) {
	Outputs := make(map[string]string)
	if lo := cfg; lo != "" {
		los := strings.Split(lo, ";")
		for _, v := range los {
			if logType2Config := strings.SplitN(v, ",", 2); len(logType2Config) == 2 {
				Outputs[logType2Config[0]] = logType2Config[1]
			} else {
				continue
			}
		}
	}
	//init log
	logs.Reset()
	for adaptor, cfg := range Outputs {
		err := logs.SetLogger(adaptor, cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s with the config %q got err:%s", adaptor, cfg, err.Error()))
		}
	}
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
}