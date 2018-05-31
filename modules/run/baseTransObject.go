package run

type BaseTransObject struct {
	Order_id      string `json:"order_id,omitempty"`      //商户订单号
	Tran_order_id string `json:"tran_order_id,omitempty"` //交易流水号        每个交易标识唯一流水号，每个子交易有各自的唯一键值
	Sys_order_id  string `json:"sys_order_id,omitempty"`  //业务流水号         一个产品对应一个业务流水号
	Acct_order_id string `json:"acct_order_id,omitempty"` //账务流水号         当交易涉及账务处理时，会对应一个账务流水号   如卡券
	Prod_cd       string `json:"prod_cd,omitempty"`
	Tran_cd       string `json:"tran_cd,omitempty"`
	Biz_cd        string `json:"biz_cd,omitempty"`
}

func (t BaseTransObject) GetOrderId() string {
	return t.Order_id
}

func (t BaseTransObject) GetTranOrderId() string {
	return t.Tran_order_id
}

func (t BaseTransObject) GetSysOrderId() string {
	return t.Sys_order_id
}

func (t BaseTransObject) GetAcctOrderId() string {
	return t.Acct_order_id
}

func (t BaseTransObject) GetProdCd() string {
	return t.Prod_cd
}

func (t BaseTransObject) GetTranCd() string {
	return t.Tran_cd
}

func (t BaseTransObject) GetBizCd() string {
	return t.Biz_cd
}
