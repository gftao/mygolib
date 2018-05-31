package channelManager

import "mygolib/gerror"

type IChainManager interface {
	RegisterChain(chainName string) gerror.IError //注册通道
	ReleaseChain(chainName string)                //释放通道
	SendMsgToChain(srcName, chainName string, msg interface{}) gerror.IError
	RecvMsgFromChain(chainName string) (string, interface{}, gerror.IError)
	CanCluster() bool //判断是否支持集群
}

var instance IChainManager

func RegisterChain(chainName string) gerror.IError {
	return instance.RegisterChain(chainName)
}
func ReleaseChain(chainName string) {
	instance.ReleaseChain(chainName)
}
func SendMsgToChain(srcName, chainName string, msg interface{}) gerror.IError {
	return instance.SendMsgToChain(srcName, chainName, msg)
}
func RecvMsgFromChain(chainName string) (string, interface{}, gerror.IError) {
	return instance.RecvMsgFromChain(chainName)
}
func CanCluster() bool {
	return instance.CanCluster()
}

//内部收发信息结构
type ChannelMsg struct {
	FromChain string
	ToChain   string
	Msg       interface{}
}
