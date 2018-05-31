package bmq

// #cgo LDFLAGS: -L/app/gol/bmq/lib -lbmqapi
// #cgo CFLAGS: -D_LINUXES_ -I/app/gol/bmq/include
/*
#include "bmq.h"
int bmqOpen(int iMbid);
int bmqClose();
int bmqPut(int iGrpid, int iMbid, int iPrior, long lType, long lClass, char *aMsgbuf, int iMsglen);
int bmqGet(int *piGrpid, int *piMbid, int *piPrior, long *plType, long *plClass, char *aMsgbuf, int *piMsglen);
int bmqGetw(int *piGrpid, int *piMbid, int *piPrior, long *plType, long *plClass, char *aMsgbuf, int *piMsglen, int iTimeout);
*/
import "C"

import "unsafe"

import (
	"fmt"
	"golib/defs"
	e "golib/gerror"
	"os"
	"runtime"
	"strconv"
)

//消息最大报文体长度
const BMQ_MAX_BUF_LEN = 8192
const BMQ_DEF_GRP_ID = 0
const BMQ_GRP_ENV = "BMQ_GROUP_ID"
const C_RETURN_SUCC = 0

var BmqGrp = BMQ_DEF_GRP_ID

func init() {
	//初始化BMQ组信息
	lGrp, err := strconv.Atoi(os.Getenv(BMQ_GRP_ENV))
	if err != nil {
		BmqGrp = lGrp
		fmt.Println("use BMQ_GRP_ENV:", BmqGrp)
	}
}

func BmqOpen(Qid int) error {
	/*线程锁*/
	runtime.LockOSThread()
	RetCd, err := C.bmqOpen(C.int(Qid))
	if RetCd != C_RETURN_SUCC {
		return e.New(int(RetCd), err, "C.bmqOpen qid[%d] fail;", Qid)
	}
	return nil
}

func BmqClose() error {
	/*线程锁*/
	runtime.LockOSThread()
	RetCd, err := C.bmqClose()
	if RetCd != C_RETURN_SUCC {
		return e.New(int(RetCd), defs.TRN_SYS_ERROR, err, "C.bmqClose error;")
	}
	return nil
}

func BmqSnd(Qid int, MsgBuf []byte) error {
	/*线程锁*/
	runtime.LockOSThread()

	var Prior C.int = 0
	var Type C.long = 0
	var Class C.long = 0

	RetCd, err := C.bmqPut(C.int(BmqGrp), C.int(Qid), Prior, Type, Class, (*C.char)(unsafe.Pointer(&MsgBuf[0])), C.int(len(MsgBuf)))
	if RetCd != C_RETURN_SUCC {
		return e.New(int(RetCd), defs.TRN_SYS_ERROR, err, "C.bmqPut qid[%d] fail;", Qid)
	}
	return nil
}

func BmqRcv() (Qid int, MsgBuf []byte, err error) {
	var Grp C.int = 0
	var Prior C.int = 0
	var Type C.long = 0
	var Class C.long = 0
	var LQid C.int = 0
	var Len C.int = BMQ_MAX_BUF_LEN

	MsgBuf = make([]byte, BMQ_MAX_BUF_LEN)

	RetCd, err := C.bmqGetw(&Grp, &LQid, &Prior, &Type, &Class, (*C.char)(unsafe.Pointer(&MsgBuf[0])), &Len, 0)
	if RetCd != C_RETURN_SUCC {
		return 0, nil, e.New(int(RetCd), defs.TRN_SYS_ERROR, err, "bmqGetW error")
	}
	return int(LQid), MsgBuf[:int(Len)], nil
}
