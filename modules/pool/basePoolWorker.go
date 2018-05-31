package pool

import "mygolib/modules/run"

/**
  基本的工作线程类  需要自己定义run方法
*/
type BasePoolWorker struct {
	run.BaseWorker
	WorkChain chan interface{}
	WorkPool  IPoolManager
}

/************ for IPoolWorker***********************/
func (t BasePoolWorker) GetId() uint32 {
	return t.Id
}
func (t BasePoolWorker) GetName() string {
	return t.NodeName
}
func (this *BasePoolWorker) Init() error {
	return nil
}

func (this *BasePoolWorker) Run() error {
	return nil
}

func (this *BasePoolWorker) SetWorkChain(ch chan interface{}) bool {
	this.WorkChain = ch
	return true
}

func (this *BasePoolWorker) SetWorkPool(pl IPoolManager) {
	this.WorkPool = pl
}

func (this *BasePoolWorker) SendToWorkChain(msg interface{}) {
	this.WorkChain <- msg
}
