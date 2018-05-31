package pool

type IPool interface {
	Take() (IPoolWorker, error)
	Return(entity IPoolWorker) error
	GetTotal() uint32
	GetUsed() uint32
}
