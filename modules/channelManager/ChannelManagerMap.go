package channelManager

import (
	"mygolib/defs"
	"mygolib/gerror"
 	"sync"
 )

type ChainManagerMap struct {
	clusterId int
	mapSize   int
	cacheSize int
	ChainMap  map[string]chan ChannelMsg
	mutex     sync.RWMutex
}

func NewChainManagerMap(mapSize, cacheSize int) IChainManager {
	return &ChainManagerMap{clusterId: 0, mapSize: mapSize, cacheSize: cacheSize,
		ChainMap: make(map[string]chan ChannelMsg, mapSize)}
}

/**************for IChainManager*****************/
func (this *ChainManagerMap) RegisterChain(chainName string) gerror.IError {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.ChainMap[chainName]; ok {
		return gerror.NewR(99001, nil, "通道名已经存在"+chainName)
	}
	c := make(chan ChannelMsg, this.cacheSize)
	this.ChainMap[chainName] = c
	return nil
}

func (this *ChainManagerMap) ReleaseChain(chainName string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.ChainMap[chainName] = nil
}

func (this *ChainManagerMap) SendMsgToChain(srcName, chainName string, msg interface{}) gerror.IError {
	if this.ChainMap[chainName] == nil {
 		return gerror.New(-1, defs.TRN_SYS_ERROR, nil, "通道名不存在，信息发送失败"+chainName)
	}
	lmsg := ChannelMsg{FromChain: srcName, ToChain: chainName, Msg: msg}
	this.ChainMap[chainName] <- lmsg
	return nil
}

func (this *ChainManagerMap) RecvMsgFromChain(chainName string) (string, interface{}, gerror.IError) {
	ch, ok := this.ChainMap[chainName]
	if !ok {
		return "", nil, gerror.New(-1, defs.TRN_SYS_ERROR, nil, "通道名不存在"+chainName)
	}
	lmsg := <-ch
	return lmsg.FromChain, lmsg.Msg, nil
}

func (this ChainManagerMap) CanCluster() bool {
	return this.clusterId > 0
}
