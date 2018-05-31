package pool

type IPoolWorker interface {
	IPoolBase
	SetWorkChain(ch chan interface{}) bool
	SetWorkPool(pl IPoolManager)
	SendToWorkChain(msg interface{})
}
