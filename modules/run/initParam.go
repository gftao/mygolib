package run

import (
	"mygolib/modules/idManager"
)

type InitParams struct {
	CluseterId      int
	SysIdGenerator  idManager.IdGenerator
	TranIdGenerator idManager.IdGenerator
}
