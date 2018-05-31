package hsm

import (
	"golib/gerror"
	"golib/modules/logs"
	"io"
	"net"
)

//全局加密机对象
var Hsm05Svr *HsmServer

func (svr *HsmServer) Init() error {

	logs.Debug("Hsm05Svr[%+v];", svr)
	for i := 0; i < svr.ConnNum; i++ {
		err := svr.Connect()
		if err != nil {
			logs.Error("加密机[%s] 连接失败[%s]", svr.HsmAddr, err)
			return err
		}
	}
	/*全部连接成功后，加密机启动正常*/
	if len(svr.ConnChan) == svr.ConnNum {
		svr.HsmStat = true
	}
	logs.Debug("Hsm05Svr 连接加密机成功[%+v] 连接数[%d];", svr, len(svr.ConnChan))
	return nil
}

func (svr *HsmServer) Connect() error {
	var err error

	peerAddr, err := net.ResolveTCPAddr("tcp4", svr.HsmAddr)
	if err != nil {
		logs.Error("ResolveTCPAddr 非法值[%s][%s];", svr.HsmAddr, err)
		return gerror.NewR(1007, err, "ResolveTCPAddr 非法值[%s]", svr.HsmAddr)
	}
	conn, err := net.DialTCP("tcp4", nil, peerAddr)
	if err != nil {
		logs.Error("net.DialTCP 失败[%s][%s];", svr.HsmAddr, err)
		return gerror.NewR(1007, err, "net.DialTCP 失败[%s]", svr.HsmAddr)
	}
	logs.Info("[%s]加密机连接成功;", svr.HsmAddr)
	conn.SetKeepAlive(true)
	svr.ConnChan <- conn
	return nil
}

func (svr *HsmServer) Close() {
	logs.Info("Hsm05Svr Begin Close;")
	for conn := range svr.ConnChan {
		conn.Close()
		if len(svr.ConnChan) == 0 {
			break
		}
	}
	logs.Info("Hsm05Svr Close Finish;")
}

func (svr *HsmServer) Comm(req *HsmRequest) (*HsmResponse, error) {
	var err error
	sndNum := 0
	reqBody := req.Byte()
	totalNum := len(reqBody)

	if svr.ConnChan == nil {
		return nil, gerror.NewR(1007, err, "加密机未初始化")
	}
	//取可用通道
	conn := <-svr.ConnChan
	for i := 0; i < totalNum; i += sndNum {
		sndNum, err = conn.Write(reqBody[i:])
		if err != nil {
			logs.Error("发送报文失败[%s]", err)
			return nil, gerror.NewR(1007, err, "发送通信报文失败:[%s];", err)
		}
		logs.Debug("发送报文成功,长度[%d][% 02X]", sndNum, reqBody[i:i+sndNum])
	}

	respBody := make([]byte, 8192)
	rcvNum, err := conn.Read(respBody)
	if err != nil && err != io.EOF {
		logs.Error("读取应答报文失败[%s];", err)
		conn.Close()
		return nil, gerror.NewR(1007, err, "读取报文失败")
	}
	logs.Debug("收到应答报文[% 02X]", respBody[:rcvNum])
	//应答通道
	svr.ConnChan <- conn
	hsmResp, err := NewHsmResponse(respBody[:rcvNum])
	if err != nil {
		logs.Error("创建NewHsmResponse对象失败[%s];", err)
		return nil, gerror.NewR(1007, err, "创建NewHsmResponse对象失败")
	}

	return hsmResp, nil
}

/*加密机请求报文*/
type HsmCall05 struct {
	inKeyIndex   int    //传入主秘钥索引
	inTekIndex   int    //TEK传输根密钥索引
	inZmk        []byte //传入区域主密钥
	inTmk        []byte //传入终端主密钥
	inKeyType    byte   //传入密钥类型
	inKeyLen     int    //传入工作密钥长度
	inWorkKey    []byte //传入密钥
	inChkValue   []byte //传入密钥校验值
	inMacAlgo    byte   //传入MAC算法
	inMac        []byte //传入MAC值，校验时使用
	inMacBlock   []byte //传入MacBlock 计算校验mac时使用
	inQPinType   byte   //传入请求Pin格式
	inQPik       []byte //传入源Pik
	inQAcct      []byte //传入源账号
	inQPinBlock  []byte //传入PinBlock
	inRPinType   byte   //传入目的Pin格式
	inRPik       []byte //传入目的Pik
	inRAcct      []byte //传入目的账号
	inRand1      []byte //传入主秘钥随机数1
	inRand2      []byte //传入主秘钥随机数2
	inReserved   []byte //传入保留参数
	outPinBlock  []byte //转换后Pin
	outTmk       []byte //返回Tmk结果
	outMac       []byte //返回Mac结果
	outWorkKey   []byte //返回工作密钥结果
	outChkValule []byte //返回工作密钥校验值
	outReserved  []byte //返回保留
	result       bool   //加密机调用结果
}

//设置计算密钥时 密钥调用
func NewHsmCall() *HsmCall05 {
	if Hsm05Svr.HsmStat != true {
		return nil
	}
	return new(HsmCall05)
}

//设置计算密钥时 密钥索引
func (opr *HsmCall05) SetInKeyIndex(index int) {
	opr.inKeyIndex = index
}

//设置计算密钥时 TEK根密钥索引
func (opr *HsmCall05) SetInTekIndex(index int) {
	opr.inTekIndex = index
}

//设置计算密钥时 通信主密钥
func (opr *HsmCall05) SetInZmk(zmk []byte) {
	opr.inZmk = zmk
}

//设置计算密钥时 终端主秘钥
func (opr *HsmCall05) SetInTmk(tmk []byte) {
	opr.inTmk = tmk
}

//设置计算密钥时 密钥长度
func (opr *HsmCall05) SetInKeyLen(len int) {
	opr.inKeyLen = len
}

//设置计算密钥时 密钥类型
func (opr *HsmCall05) SetInKeyType(tp byte) {
	if tp != WK_COM_TP &&
		tp != WK_PIK_TP &&
		tp != WK_DAT_TP &&
		tp != WK_MAK_TP {
		logs.Error("SetInKeyType type[%c] 非法;", tp)
		return
	}
	opr.inKeyType = tp
}

//设置计算密钥时 工作密钥
func (opr *HsmCall05) SetWorkKey(key []byte) {
	opr.inWorkKey = key
}

//设置计算密钥时 校验值
func (opr *HsmCall05) SetInChkValue(chkVal []byte) {
	opr.inChkValue = chkVal
}

//设置转PIN时 源Pin类型
func (opr *HsmCall05) SetInQPinType(tp byte) {
	opr.inQPinType = tp
}

//设置转PIN时 源Pik
func (opr *HsmCall05) SetInQPik(pik []byte) {
	opr.inQPik = pik
}

//设置转PIN时 源账号
func (opr *HsmCall05) SetInQAcct(acct []byte) {
	opr.inQAcct = acct
}

//设置转PIN时 源PinBlock
func (opr *HsmCall05) SetInQPinBlock(block []byte) {
	opr.inQPinBlock = block
}

//设置转PIN时 目的Pin类型
func (opr *HsmCall05) SetInRPinType(tp byte) {
	opr.inRPinType = tp
}

//设置转PIN时 目的Pik
func (opr *HsmCall05) SetInRPik(pik []byte) {
	opr.inRPik = pik
}

//设置转PIN时 目的异或账号
func (opr *HsmCall05) SetInRAcct(acct []byte) {
	opr.inRAcct = acct
}

//设置Mac计算算法
func (opr *HsmCall05) SetInMacAlgo(algo byte) {
	opr.inMacAlgo = algo
}

// MAC 计算时传入MAB
func (opr *HsmCall05) SetInMacBlock(block []byte) {
	opr.inMacBlock = block
}

// MAC 校验时传入目的Mac
func (opr *HsmCall05) SetInMac(mac []byte) {
	opr.inMac = mac
}

// 生成主密钥时随机因子1
func (opr *HsmCall05) SetInRand1(rand []byte) {
	opr.inRand1 = rand
}

// 生成主密钥时随机因子2
func (opr *HsmCall05) SetInRand2(rand []byte) {
	opr.inRand2 = rand
}

//设置保留域
func (opr *HsmCall05) SetInReserved(resv []byte) {
	opr.inReserved = resv
}

//设置调用结果(内部)
func (opr *HsmCall05) setCallResult(r bool) {
	opr.result = r
}

//返回调用结果
func (opr *HsmCall05) GetCallResult() bool {
	return opr.result
}

//返回终端主秘钥
func (opr *HsmCall05) GetTmk() []byte {
	return opr.outTmk
}

//返回工作密钥
func (opr *HsmCall05) GetWorkKey() []byte {
	return opr.outWorkKey
}

//返回校验值
func (opr *HsmCall05) GetCheckValue() []byte {
	return opr.outChkValule
}

//返回PinBlock
func (opr *HsmCall05) GetPinBlock() []byte {
	return opr.outPinBlock
}

//返回Mac
func (opr *HsmCall05) GetMac() []byte {
	return opr.outMac
}

//返回保留域
func (opr *HsmCall05) GetReserved() []byte {
	return opr.outReserved
}

/*
必输: Mac算法 Mac工作密钥 MacBlock
响应：GetMac
*/
func (opr *HsmCall05) GenMac() error {
	var err error
	Req := NewHsmRequest()
	Req.SetAsc("D132")
	Req.SetByte(opr.inMacAlgo)
	Req.SetInt8(len(opr.inWorkKey) / 2)
	Req.SetAsc(string(opr.inWorkKey))
	Req.SetHex([]byte(IV))
	Req.SetInt16(len(opr.inMacBlock))
	Req.SetHex(opr.inMacBlock)

	logs.Debug("GenMac 请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if Rsp.CheckSucc() {
		out, err := Rsp.GetAsc(8)
		if err != nil {
			opr.setCallResult(false)
			logs.Error("获取MAC 结果失败:[%s]", err)
			return gerror.NewR(9004, err, "获取MAC结果失败;")
		}
		logs.Debug("计算MAC结果:[%s]", out)
		opr.outMac = []byte(out)
		opr.setCallResult(true)
		return nil
	} else {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("计算MAC失败:[%0X]", emsg)
		return gerror.NewR(9004, err, "计算MAC失败;")
	}
}

/*
必输： Mac 算法； 工作密钥; Mac; MacBlock;
得到结果: GetCallResult()
*/
func (opr *HsmCall05) VerifyMac() error {
	var err error
	Req := NewHsmRequest()
	Req.SetAsc("D134")
	Req.SetByte(opr.inMacAlgo) // Mac 算法
	Req.SetInt8(len(opr.inWorkKey) / 2)
	Req.SetAsc(string(opr.inWorkKey)) // 工作密钥
	Req.SetHex([]byte(IV))
	Req.SetAsc(string(opr.inMac))
	Req.SetInt16(len(opr.inMacBlock))
	Req.SetHex([]byte(opr.inMacBlock))

	logs.Debug("VerifyMac 请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if Rsp.CheckSucc() {
		out, err := Rsp.GetAsc(8)
		if err != nil {
			opr.setCallResult(false)
			logs.Error("VerifyMac 失败:[%s]", err)
			return gerror.NewR(9004, err, "VerifyMac MAC 失败;")
		}
		logs.Debug("VerifyMac MAC结果:[%s]", out)
		opr.outMac = []byte(out)
		opr.setCallResult(true)
		return nil
	} else {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("计算MAC失败:[%0X]", emsg)
		return gerror.NewR(9004, err, "计算MAC失败;")
	}
	return nil
}

/*
必输：源Pik 源PinType 源PinBlock 源PinAcct
      目的Pik  目的PinType 目的PinAcct

*/
func (opr *HsmCall05) ConvPin() error {
	var err error
	Req := NewHsmRequest()
	Req.SetAsc("D124")
	Req.SetInt8(len(opr.inQPik) / 2)    //源Pik 长度
	Req.SetAsc(string(opr.inQPik))      //源Pik
	Req.SetInt8(len(opr.inRPik) / 2)    //目的Pik 长度
	Req.SetAsc(string(opr.inRPik))      //目的Pik
	Req.SetByte(opr.inQPinType)         //源Pin 类型
	Req.SetByte(opr.inRPinType)         //目的Pin 类型
	Req.SetAsc(string(opr.inQPinBlock)) //源 PinBlock
	Req.SetHex(opr.inQAcct)             //源 账号
	Req.SetByte(';')
	Req.SetHex(opr.inRAcct) //目的 账号
	Req.SetByte(';')

	logs.Debug("ConvPin 请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if Rsp.CheckSucc() {
		out, err := Rsp.GetAsc(8)
		if err != nil {
			opr.setCallResult(false)
			logs.Error("ConvPin 失败:[%s]", err)
			return gerror.NewR(9004, err, "ConvPin 失败;")
		}
		logs.Debug("ConvPin 结果:[%s]", out)
		opr.outPinBlock = []byte(out)
		opr.setCallResult(true)
		return nil
	} else {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("ConvPin 失败:[%0X]", emsg)
		return gerror.NewR(9004, err, "ConvPin 失败:;")
	}
	return nil
}

/*导入工作密钥:必输: 索引/主秘钥  密钥类型  工作密钥

响应：GetWorkKey()
	  GetCheckVal()
*/
func (opr *HsmCall05) ImportWorkKey() error {
	var err error
	Req := NewHsmRequest()
	Req.SetAsc("D102")
	Req.SetInt8(len(opr.inWorkKey) / 2) //工作密钥长度
	Req.SetInt8(16)
	if opr.inKeyIndex > 0 {
		Req.SetInt16(opr.inKeyIndex) //通信主密钥索引
	} else {
		Req.SetInt16(0xffff)
	}
	Req.SetByte(opr.inKeyType) //工作密钥类型
	if opr.inKeyIndex == 0 {
		Req.SetAsc(string(opr.inZmk)) //通信主秘钥
	}
	Req.SetAsc(string(opr.inWorkKey)) //工作密钥

	if len(opr.inChkValue) > 0 {
		Req.SetInt8(len(opr.inChkValue)/2)
		Req.SetAsc(string(opr.inChkValue))
	}

	logs.Debug("ImportWorkKey 请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if Rsp.CheckSucc() {
		len, err := Rsp.GetInt8()
		if err != nil {
			opr.setCallResult(false)
			logs.Error("获取返回密钥长度失败:[%s]", err)
			return gerror.NewR(9004, err, "获取返回密钥长度失败;")
		}
		out, err := Rsp.GetAsc(len)
		if err != nil {
			opr.setCallResult(false)
			logs.Error("获取返回密钥失败:[%s]", err)
			return gerror.NewR(9004, err, "获取返回密钥失败;")
		}
		logs.Debug("返回工作密钥:[%s]", out)
		opr.outWorkKey = []byte(out)
		chk, err := Rsp.GetAsc(8)
		if err != nil {
			opr.setCallResult(false)
			logs.Error("获取返回校验值失败:[%s]", err)
			return gerror.NewR(9004, err, "获取返回校验值失败;")
		}
		opr.outChkValule = []byte(chk)
		opr.setCallResult(true)
		return nil
	} else {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("密钥导入失败:[%0X]", emsg)
		return gerror.NewR(9004, err, "密钥导入失败;")
	}
	return nil
}

/*
  必输： 密钥长度  密钥索引 通信主密钥 工作密钥
  响应： GetWorkKey()
         GetCheckVal()
*/
func (opr *HsmCall05) ExportWorkKey() error {
	var err error
	Req := NewHsmRequest()
	Req.SetAsc("D104")
	Req.SetInt8(len(opr.inWorkKey) / 2) //工作密钥长度
	Req.SetInt8(16)
	if opr.inKeyIndex > 0 {
		Req.SetInt16(opr.inKeyIndex) //通信主密钥索引
	} else {
		Req.SetInt16(0xffff)
	}
	Req.SetByte(opr.inKeyType) //工作密钥类型
	if opr.inKeyIndex == 0 {
		Req.SetAsc(string(opr.inZmk)) //通信主秘钥
	}
	Req.SetAsc(string(opr.inWorkKey)) //工作密钥

	if len(opr.inChkValue) > 0 {
		Req.SetInt8(len(opr.inChkValue))
		Req.SetAsc(string(opr.inChkValue))
	}

	logs.Debug("ExportWorkKey 请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if Rsp.CheckSucc() {
		len, err := Rsp.GetInt8()
		if err != nil {
			opr.setCallResult(false)
			logs.Error("获取返回密钥长度失败:[%s]", err)
			return gerror.NewR(9004, err, "获取返回密钥长度失败;")
		}
		out, err := Rsp.GetAsc(len)
		logs.Debug("返回工作密钥:[%s]", out)
		opr.outWorkKey = []byte(out)
		chk, err := Rsp.GetAsc(8)
		if err != nil {
			opr.setCallResult(false)
			logs.Error("获取返回密钥校验值失败:[%s]", err)
			return gerror.NewR(9004, err, "获取返回密钥校验值失败;")
		}
		opr.outChkValule = []byte(chk)
		opr.setCallResult(true)
		return nil
	} else {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("密钥导出失败:[%0X]", emsg)
		return gerror.NewR(9004, err, "密钥导入失败;")
	}
	return nil
}

/*
必输:TMK根索引
	 TEK根索引
	 SEK明文（ZMK）
	 Rand1 发散因子 （ASCII）
	 Rand2 发散银子 （ASCII）
应答：
	GetTmk()
	GetChkValue()
*/
//设置计算密钥时 密钥类型

func (opr *HsmCall05) GenTermTmk() error {
	var err error

	/*TMK发散生成*/
	Req := NewHsmRequest()
	Req.SetAsc("D181")
	Req.SetInt16(opr.inKeyIndex)  //TMK 根索引
	Req.SetInt16(opr.inTekIndex)  //TEK 根索引
	Req.SetByte(KEY_TMK_TEK_TP)   //生成终端主密钥
	Req.SetInt8(len(opr.inRand1)) //发散因子1长度
	Req.SetHex(opr.inRand1)       //发散因子1
	Req.SetInt8(len(opr.inRand2)) //发散因子2长度
	Req.SetHex(opr.inRand2)       //发散因子2

	logs.Debug("D181 发散终端主密钥请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if !Rsp.CheckSucc() {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("发散终端主密钥请求:[%0X]", emsg)
		return gerror.NewR(9004, err, "密钥导入失败;")
	}

	tmkLen, err := Rsp.GetInt8()
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回密钥长度失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回密钥长度失败;")
	}

	lmkTmk, err := Rsp.GetAsc(tmkLen)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回终端主密钥失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回终端主密钥失败;")
	}
	logs.Debug("返回终端主密钥:[%s]", lmkTmk)
	_, err = Rsp.GetHex(tmkLen)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回终端主密钥TMK失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回终端主密钥TMK失败;")
	}
	chkVal, err := Rsp.GetAsc(8)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回终端主密钥TMK失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回终端主密钥TEK失败;")
	}
	logs.Debug("返回终端主密钥CheckVal:[%s]", chkVal)

	/*TEK导入*/
	Req = NewHsmRequest()
	Req.SetAsc("D108")
	Req.SetInt8(len(opr.inZmk) / 2) //传入TEK明文
	Req.SetByte(WK_COM_TP)          //通信主密钥
	Req.SetAsc(string(opr.inZmk))   //TEK明文
	logs.Debug("D108 转换终端TEK密钥请求[%s]", Req)
	Rsp, err = Hsm05Svr.Comm(Req)
	if err != nil {
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if !Rsp.CheckSucc() {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("转换终端TEK密钥失败:[%0X]", emsg)
		return gerror.NewR(9004, err, "密钥导入失败;")
	}

	tekLen, err := Rsp.GetInt8()
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回密钥长度失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回密钥长度失败;")
	}
	lmkTek, err := Rsp.GetAsc(tekLen)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回终端主密钥失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回终端主密钥失败;")
	}
	logs.Debug("转换终端TEK密钥:[%s]", lmkTek)

	/*TMK转换*/
	lopr := NewHsmCall()
	lopr.SetInKeyType(WK_COM_TP)
	lopr.SetInZmk([]byte(lmkTek))
	lopr.SetWorkKey([]byte(lmkTmk))
	lopr.SetInChkValue([]byte(chkVal))
	err = lopr.ExportWorkKey()
	if err != nil {
		logs.Error("lopr.ExportWorkKey转换终端主密钥失败:[%s]", err)
		opr.setCallResult(false)
		return gerror.NewR(9004, err, "ExportWorkKey转换终端主密钥失败;")
	}

	opr.outTmk = lopr.GetWorkKey()
	opr.outChkValule = lopr.GetCheckValue()
	logs.Debug("生成终端主密钥TMK[%s] CHECKVAL[%s];", opr.outTmk, opr.outChkValule)
	opr.setCallResult(true)
	return nil
}

/*
	必输：密钥类型  工作密钥长度  发散因子1  发散因子2
	应答：
	GetTmk()
	GetWorkKey()
	GetChkValue()
*/
func (opr *HsmCall05) GenTermWorkKey() error {
	var err error

	/*TMK发散生成*/
	Req := NewHsmRequest()
	Req.SetAsc("D182")
	Req.SetInt16(opr.inKeyIndex)  //根密钥索引
	Req.SetByte(opr.inKeyType)    //密钥类型
	Req.SetInt8(opr.inKeyLen)     //长度
	Req.SetInt8(len(opr.inRand1)) //发散因子1长度
	Req.SetHex(opr.inRand1)       //发散因子1
	Req.SetInt8(len(opr.inRand2)) //发散因子2长度
	Req.SetHex(opr.inRand2)       //发散因子2

	logs.Debug("D182 发散终端主密钥请求[%s]", Req)
	Rsp, err := Hsm05Svr.Comm(Req)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("加密机Com失败[%s]", err)
		return gerror.NewR(9003, err, "加密机Com失败;")
	}
	if !Rsp.CheckSucc() {
		opr.setCallResult(false)
		emsg := Rsp.GetLByte()
		logs.Error("发散终端主密钥请求:[%0X]", emsg)
		return gerror.NewR(9004, err, "密钥导入失败;")
	}

	tmkLen, err := Rsp.GetInt8()
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回密钥长度失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回密钥长度失败;")
	}

	zmkTmk, err := Rsp.GetAsc(tmkLen)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("返回ZMK工作密钥:[%s]", err)
		return gerror.NewR(9004, err, "返回ZMK工作密钥;")
	}
	logs.Debug("返回ZMK工作密钥:[%s]", zmkTmk)
	lmkTmk, err := Rsp.GetAsc(tmkLen)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("返回LMK工作密钥:[%s]", err)
		return gerror.NewR(9004, err, "返回LMK工作密钥;")
	}
	logs.Debug("返回LMK工作密钥:[%s]", lmkTmk)
	chkVal, err := Rsp.GetAsc(8)
	if err != nil {
		opr.setCallResult(false)
		logs.Error("获取返回终端工作密钥chkVal失败:[%s]", err)
		return gerror.NewR(9004, err, "获取返回终端工作密钥chkVal失败;")
	}
	logs.Debug("返回LMK终端工作密钥校验值:[%s]", chkVal)
	opr.outWorkKey = []byte(zmkTmk)
	opr.outTmk = []byte(lmkTmk)
	opr.outChkValule = []byte(chkVal)

	opr.setCallResult(true)
	return nil
}
