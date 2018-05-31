package hsm

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golib/gerror"
	"golib/modules/logs"
	"strings"
)

/*请求报文*/
type HsmRequest struct {
	buf []byte
	len int
}

func NewHsmRequest() *HsmRequest {
	req := new(HsmRequest)
	req.buf = make([]byte, MaxReqLen)
	req.len = 0
	return req
}

func (req *HsmRequest) SetByte(b byte) {
	req.buf[req.len] = b
	req.len++
}

func (req *HsmRequest) SetInt8(i int) {
	req.buf[req.len] = byte(i)
	req.len += 1
}

func (req *HsmRequest) SetInt16(i int) {
	binary.BigEndian.PutUint16(req.buf[req.len:], uint16(i))
	req.len += 2
}

func (req *HsmRequest) SetHex(hex []byte) {
	copy(req.buf[req.len:], hex)
	req.len += len(hex)
}

func (req *HsmRequest) SetAsc(asc string) {
	bt, _ := hex.DecodeString(asc)
	copy(req.buf[req.len:], bt)
	req.len += len(bt)

}

func (req *HsmRequest) String() string {
	req.Byte()
	return fmt.Sprintf("[%04d] [% 02X]", req.len, req.buf[:req.len])
}
func (req *HsmRequest) Byte() []byte {
	return req.buf[:req.len:req.len]
}

type HsmResponse struct {
	respBody []byte
	retCd    byte
	len      int
	currBuf  []byte
	currLen  int
}

func NewHsmResponse(buf []byte) (*HsmResponse, error) {
	if len(buf) < 1 {
		logs.Error("应答报文为空;")
		return nil, gerror.NewR(9000, nil, "应答报文为空")
	}
	hsmResp := new(HsmResponse)
	hsmResp.respBody = buf
	hsmResp.len = len(buf)

	hsmResp.retCd = hsmResp.respBody[0]
	hsmResp.currBuf = hsmResp.respBody[1:]
	hsmResp.currLen = hsmResp.len - 1

	return hsmResp, nil
}

func (rsp *HsmResponse) GetRetCd() byte {
	return rsp.retCd
}

func (rsp *HsmResponse) GetByte() (byte, error) {
	if rsp.currLen >= 1 {
		bt := rsp.currBuf[0]
		rsp.currBuf = rsp.currBuf[1:]
		rsp.currLen--
		return bt, nil
	} else {
		logs.Error("报文解析完毕，剩余为空;")
		return 0, gerror.NewR(9001, nil, "报文解析完毕，剩余为空;")
	}
}

func (rsp *HsmResponse) GetInt16() (uint16, error) {
	if rsp.currLen >= 2 {
		i := binary.BigEndian.Uint16(rsp.currBuf[:2])
		rsp.currBuf = rsp.currBuf[2:]
		rsp.currLen -= 2
		return i, nil
	} else {
		logs.Error("报文解析完毕，剩余为空;")
		return 0, gerror.NewR(9002, nil, "报文解析完毕，剩余为空;")
	}
}

func (rsp *HsmResponse) GetInt8() (int, error) {
	if rsp.currLen >= 1 {
		i := int(rsp.currBuf[0])
		rsp.currBuf = rsp.currBuf[1:]
		rsp.currLen -= 1
		return i, nil
	} else {
		logs.Error("报文解析完毕，剩余为空;")
		return 0, gerror.NewR(9002, nil, "报文解析完毕，剩余为空;")
	}
}

func (rsp *HsmResponse) GetHex(len int) ([]byte, error) {
	if rsp.currLen >= len {
		hex := rsp.currBuf[:len]
		rsp.currBuf = rsp.currBuf[len:]
		rsp.currLen -= len
		return hex, nil
	} else {
		logs.Error("报文解析完毕，剩余为空;")
		return nil, gerror.NewR(9002, nil, "报文解析完毕，剩余为空;")
	}
}

func (rsp *HsmResponse) GetAsc(len int) (string, error) {
	if rsp.currLen >= len {
		bin := rsp.currBuf[:len]
		rsp.currBuf = rsp.currBuf[len:]
		rsp.currLen -= len
		asc := strings.ToUpper(hex.EncodeToString(bin))
		return asc, nil
	} else {
		logs.Error("报文解析完毕，剩余为空;")
		return "", gerror.NewR(9002, nil, "报文解析完毕，剩余为空;")
	}
}

func (rsp *HsmResponse) GetLByte() []byte {
	return rsp.currBuf
}

func (rsp *HsmResponse) CheckSucc() bool {
	return rsp.retCd == success
}

func (rsp *HsmResponse) String() string {
	return fmt.Sprintf("应答[% 02X][%c][%d][% 02X]",
		rsp.respBody, rsp.retCd, rsp.currLen, rsp.currBuf)
}
