package run

import (
	"mygolib/gerror"
	"net/http"
)

type ICommOperator interface {
	UnPackMsg(r *http.Request) gerror.IError //解析报文
	VerifyMsg() gerror.IError                //校验报文
	SignMsg() gerror.IError                  //签名
	PackMsg() gerror.IError                  //组报文
}
