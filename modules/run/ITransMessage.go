package run

type ITransMsg interface {
	GetOrderId() string
	GetTranOrderId() string
	GetSysOrderId() string
	GetAcctOrderId() string
	GetProdCd() string
	GetTranCd() string
	GetBizCd() string
	ToString() string
}
