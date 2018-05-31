package cache

import (
	"errors"
	"mygolib/modules/config"
 	"time"
	"mygolib/modules/myLogger"
)

var initFlg bool = false

func InitModule() error {

	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}

	clusterId := config.IntDefault("clusterId", 0)

	cacheType := config.StringDefault("cacheType", "Mem")
	expTime := config.IntDefault("expTime", 5)

	ok := initCache(clusterId, cacheType, time.Duration(expTime))

	if !ok {
		return errors.New("初始化cache失败")
	}

	myLogger.Info("初始化cache成功")
	initFlg = true

	return nil
}

func HasModuleInit() bool {
	return initFlg
}

func initCache(clusterId int, cacheType string, expTime time.Duration) bool {

	// Set the default expiration time.
	if expTime == 0 {
		expTime = time.Minute
	}

	switch cacheType {
	case "Mem":
		if clusterId > 0 {
			myLogger.Error("集群模式不能使用MEM缓存")
			return false
		}
		Instance = NewInMemoryCache(expTime)
	default:
		myLogger.Error("非法的缓存类别", cacheType)
		return false
	}

	myLogger.Info("创建缓存成功")
	return true
}
