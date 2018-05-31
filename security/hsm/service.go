package hsm

import (
	"golib/gerror"
	"golib/modules/logs"
)

const (
	//MAC向量
	IV = "\x00\x00\x00\x00\x00\x00\x00\x00"
	//MAC算法
	MAC_ALGO_XOR = 0x01
	MAC_ALGO_99  = 0x02
	MAC_ALGO_919 = 0x03
	MAC_ALGO_ECB = 0x04
	//PIN格式
	PIN_BLOCK_01 = 0x01
	PIN_BLOCK_02 = 0x02
	PIN_BLOCK_03 = 0x03
	PIN_BLOCK_04 = 0x04
	PIN_BLOCK_05 = 0x05
	PIN_BLOCK_06 = 0x06
	//指令类型
	CMD_ENC = 0x00
	CMD_DEC = 0x01
	//生成终端主秘钥类型 D181使用
	KEY_TEK_TP     = 0x00
	KEY_TMK_TP     = 0x01
	KEY_TMK_TEK_TP = 0x02
	//工作密钥类型-变种
	WK_COM_TP = 0x01
	WK_PIK_TP = 0x11
	WK_MAK_TP = 0x12
	WK_DAT_TP = 0x13

	MaxReqLen = 4096
	success   = 'A'
)

type HsmOper interface {
	/* 设置区域主密钥索引*/
	SetInKeyIndex(index int)
	/* 设置TEK主密钥索引*/
	SetInTekIndex(index int)
	/* 设置区域主密钥*/
	SetInZmk(zmk []byte)
	/* 设置终端主密钥*/
	SetInTmk(tmk []byte)
	/* 设置密钥长度*/
	SetInKeyLen(len int)
	/* 设置密钥类型*/
	SetInKeyType(tp byte)
	/* 设置密钥*/
	SetWorkKey(key []byte)
	/* 设置校验值*/
	SetInChkValue(chkVal []byte)
	SetInQPinType(tp byte)
	SetInQPik(pik []byte)
	SetInQAcct(acct []byte)
	SetInQPinBlock(block []byte)
	SetInRPinType(tp byte)
	SetInRPik(pik []byte)
	SetInRAcct(acct []byte)
	SetInMacAlgo(algo byte)
	SetInMacBlock(block []byte)
	SetInMac(mac []byte)
	SetInRand1(rand []byte)
	SetInRand2(rand []byte)
	SetInReserved(resv []byte)
	/*加密机操作*/
	ImportWorkKey() error
	/*
	   必输：
	          密钥长度  密钥索引 通信主密钥 工作密钥
	   响应： GetWorkKey()
	          GetCheckVal()
	*/
	ExportWorkKey() error
	/*
		必输:
			Mac算法 Mac工作密钥 MacBlock
		响应：
			GetMac
	*/
	GenMac() error
	/*
		必输： Mac 算法； 工作密钥; Mac; MacBlock;
		得到结果: GetCallResult()
	*/
	VerifyMac() error
	/*
		必输：源Pik 源PinType 源PinBlock 源PinAcct
		      目的Pik  目的PinType 目的PinAcct

	*/
	ConvPin() error
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
	GenTermTmk() error

	/*
		必输：
			根密钥索引 密钥类型  工作密钥长度  发散因子1  发散因子2
		应答：
			GetTmk()
			GetWorkKey()
			GetChkValue()
	*/
	GenTermWorkKey() error

	/*取应答结果*/
	GetCallResult() bool
	GetTmk() []byte
	GetWorkKey() []byte
	GetCheckValue() []byte
	GetMac() []byte
	GetPinBlock() []byte
}

func NewHsmOper() (HsmOper, error) {
	if GlbSvr == nil {
		logs.Error("加密机未初始化; 或者加密机未操作成功;")
		return nil, gerror.NewR(9999, nil, "加密机未初始化; 或者加密机未操作成功;")
	}
	switch GlbSvr.HsmType {
	case HSM_SJL05:
		if Hsm05Svr == nil {
			logs.Error("HSM_SJL05加密机未初始化; 或者加密机未操作成功;")
			return nil, gerror.NewR(9999, nil, "加密机未初始化; 或者加密机未操作成功;")
		}
		c := NewHsmCall()
		if c == nil {
			return nil, gerror.NewR(9999, nil, "加密机未初始化成功;")
		}
		return c, nil
	default:
		return nil, gerror.NewR(9999, nil, "非法加密机类型;")
	}
}
