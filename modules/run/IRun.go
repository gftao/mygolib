package run

import "mygolib/gerror"

type IRun interface {
	Init(initParams InitParams, chainName string) gerror.IError
	Run()
	Finish()
}
