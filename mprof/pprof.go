package mprof

import (
	_ "net/http/pprof"
	"golib/modules/config"
	"net/http"
)

func InitModel() {
	var ppaddr string
	if config.HasModuleInit() {
		ppaddr = config.StringDefault("ppaddr", "")
		if ppaddr != "" {
			go http.ListenAndServe(ppaddr, nil)
		}
	}

}