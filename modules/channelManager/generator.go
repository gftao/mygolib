package channelManager

import (
	"errors"
	"golib/modules/config"
	"golib/modules/logr"
)

var initFlg bool = false

func InitModule() error {

	if !config.HasModuleInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}

	clusterId := config.IntDefault("clusterId", 0)
	chainType := config.StringDefault("chainType", "Map")
	mapSize := config.IntDefault("mapSize", 10)
	cacheSize := config.IntDefault("cacheSize", 100)

	var ok bool
	instance, ok = GetChainManager(chainType, clusterId, mapSize, cacheSize)

	if !ok {
		return errors.New("初始化cache失败")
	}

	initFlg = true

	return nil
}

func HasModuleInit() bool {
	return initFlg
}

func GetChainManager(chainType string, clusterId int,
	mapSize, cacheSize int) (ch IChainManager, ok bool) {
	switch chainType {
	case "Map":
		if clusterId != 0 {
			logr.Error(clusterId, "集群模式不能使用 MAP 通道管理器")
			return nil, false
		}
		return NewChainManagerMap(mapSize, cacheSize), true
	}
	return nil, false
}
