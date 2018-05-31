package idManager

import (
	"mygolib/defs"
	"sync"
)

// ID生成器的接口类型。
type IdGenerator interface {
	GetUint32() uint32 // 获得一个uint32类型的ID。
}

// 创建ID生成器。
func NewUIdGenerator(clusterId int) IdGenerator {
	return &cyclicIdGenerator{clusterId: clusterId, sn: 0, ended: false}
}

// ID生成器的实现类型。
type cyclicIdGenerator struct {
	clusterId int
	sn        uint32     // 当前的ID。
	ended     bool       // 前一个ID是否已经为其类型所能表示的最大值。
	mutex     sync.Mutex // 互斥锁。
}

func (this *cyclicIdGenerator) GetUint32() uint32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.ended {
		defer func() { this.ended = false }()
		this.sn = 0
		return this.sn
	}
	id := this.sn
	if id < defs.CYCLEMAX {
		this.sn++
	} else {
		this.ended = true
	}
	id = id%defs.CYCLEMAX + uint32(this.clusterId)*defs.CYCLEMAX
	return id
}
