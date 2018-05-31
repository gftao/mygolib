package pool

import (
	"errors"
	"fmt"
	"mygolib/modules/channelManager"
	"mygolib/modules/run"
	"os"
)

/**
  管理池基本类
  要实现池，只需要继承该结构，并实现自定义的创建方法
*/
type BasePoolManager struct {
	run.BaseWorker
	WorkPool      IPool
	SelfPoolChain chan interface{}
}

/**************for IPoolManager*******************/
func (t BasePoolManager) GetId() uint32 {
	return t.Id
}
func (t BasePoolManager) GetName() string {
	return t.NodeName
}
func (t *BasePoolManager) Init() error {
	return errors.New("Init方法不可用")
}

func (t *BasePoolManager) Run() error {
	for i := uint32(0); i < t.WorkPool.GetTotal(); i++ {
		worker, err := t.WorkPool.Take()
		if err != nil {
			fmt.Println(t.NodeName+"启动失败", err)
			os.Exit(-1)
		}
		go worker.Run()
	}

	go func() {
		for {
			src, msg, gerr := channelManager.RecvMsgFromChain(t.NodeName)
			if gerr != nil {
				fmt.Println(t.NodeName+"读信息失败", gerr)
				os.Exit(-1)
			}
			t.Debugf("收到%s发过来的信息", src)
			t.SelfPoolChain <- msg
		}
	}()

	return nil
}
