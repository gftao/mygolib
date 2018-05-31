package gerror

import (
	"fmt"
	"mygolib/defs"
	"runtime"
	"time"
)

const (
	ErrCodeOK                 = 0
	ErrCodeInvalidCredential  = 40001
	ErrCodeAccessTokenExpired = 42001
)

type Error struct {
	ErrCode int    `json:"errcode"` //错误码
	RespCd  string `json:"respcd"`  //应答码
	ErrMsg  string `json:"errmsg"`  //错误信息
	ErrFile string `json:"errfile"` //报错源码文件名
	ErrLine int    `json:"errline"` //报错源码文件行
	ErrFunc string `json:"errfunc"` //报错函数名
	Err     error  `json:"err"`     //传入错误
	ErrTm   string `json:"errtm"`   //报错时间
}

func New(ErrCode int, RespCd string, Err error, format string, v ...interface{}) IError {

	funcName, file, line, ok := runtime.Caller(1)
	if !ok {
	}
	return &Error{
		ErrTm:   time.Now().Format("2006-01-02 15:04:05.000000"),
		ErrCode: ErrCode,
		RespCd:  RespCd,
		Err:     Err,
		ErrMsg:  fmt.Sprintf(format, v...),
		ErrFile: file,
		ErrLine: line,
		ErrFunc: runtime.FuncForPC(funcName).Name()}
}

func NewR(ErrCode int, Err error, format string, v ...interface{}) IError {

	funcName, file, line, ok := runtime.Caller(1)
	if !ok {
	}
	return &Error{
		ErrTm:   time.Now().Format("2006-01-02 15:04:05.000000"),
		ErrCode: ErrCode,
		RespCd:  defs.TRN_SYS_ERROR,
		Err:     Err,
		ErrMsg:  fmt.Sprintf(format, v...),
		ErrFile: file,
		ErrLine: line,
		ErrFunc: runtime.FuncForPC(funcName).Name()}
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%v] [%s]:[%d] [%v]-C[%d]R[%s]E[%v]M[%v];\n",
		e.ErrTm, e.ErrFile, e.ErrLine, e.ErrFunc, e.ErrCode, e.RespCd, e.Err, e.ErrMsg)
}

//for IError
func (e *Error) GetErrorCode() string {
	return e.RespCd
}

func (e *Error) GetErrorString() string {
	return e.ErrMsg
}
