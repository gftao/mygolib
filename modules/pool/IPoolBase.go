package pool

type IPoolBase interface {
	GetId() uint32
	GetName() string
	Init() error
	Run() error
}
