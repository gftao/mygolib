package pool

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Pool struct {
	Total       uint32
	Etype       reflect.Type
	GenEntity   func() IPoolWorker
	container   chan IPoolWorker
	IdContainer map[uint32]bool
	mutex       sync.Mutex
}

// 创建实体池。
func NewPool(total uint32, genEntity func() IPoolWorker) (IPool, error) {
	if total == 0 {
		errMsg :=
			fmt.Sprintf("The pool can not be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}

/*	if genEntity() == nil {
		errMsg :=
			fmt.Sprintf("The type of result of function genEntity() is nil")
		return nil, errors.New(errMsg)
	}*/

	//entityType := reflect.TypeOf(genEntity())
	size := int(total)
	container := make(chan IPoolWorker, size)
	idContainer := make(map[uint32]bool)
	for i := 0; i < size; i++ {
		newEntity := genEntity()
		if newEntity == nil {
			errMsg :=
				fmt.Sprintf("The value of result of function genEntity() is nil")
			return nil, errors.New(errMsg)
		}
		container <- newEntity
		idContainer[newEntity.GetId()] = true
	}
	pool := &Pool{
		Total:       total,
		//Etype:       entityType,
		GenEntity:   genEntity,
		container:   container,
		IdContainer: idContainer,
	}
	return pool, nil
}

/************************for IPool***************************/

func (this *Pool) Take() (IPoolWorker, error) {
	entity, ok := <-this.container
	if !ok {
		return nil, errors.New("The inner container is invalid!")
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.IdContainer[entity.GetId()] = false
	return entity, nil
}

func (this *Pool) Return(entity IPoolWorker) error {
	if entity == nil {
		return errors.New("The returning entity is invalid!")
	}
	/*if this.Etype != reflect.TypeOf(entity) {
		errMsg := fmt.Sprintf("The type of returning entity is NOT %s!\n", this.Etype)
		return errors.New(errMsg)
	}*/
	entityId := entity.GetId()
	casResult := this.compareAndSetForIdContainer(entityId, false, true)
	if casResult == 1 {
		this.container <- entity
		return nil
	} else if casResult == 0 {
		errMsg := fmt.Sprintf("The entity (id=%d) is already in the pool!\n", entityId)
		return errors.New(errMsg)
	} else {
		errMsg := fmt.Sprintf("The entity (id=%d) is illegal!\n", entityId)
		return errors.New(errMsg)
	}
}

// 比较并设置实体ID容器中与给定实体ID对应的键值对的元素值。
// 结果值：
//       -1：表示键值对不存在。
//        0：表示操作失败。
//        1：表示操作成功。
func (this *Pool) compareAndSetForIdContainer(
	entityId uint32, oldValue bool, newValue bool) int8 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	v, ok := this.IdContainer[entityId]
	if !ok {
		return -1
	}
	if v != oldValue {
		return 0
	}
	this.IdContainer[entityId] = newValue
	return 1
}

func (t Pool) GetTotal() uint32 {
	return t.Total
}

func (t Pool) GetUsed() uint32 {
	return t.Total - uint32(len(t.container))
}
