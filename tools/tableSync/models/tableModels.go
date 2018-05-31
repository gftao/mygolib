package models

import (
	"time"
	"reflect"
	"github.com/jinzhu/gorm"
	"golib/modules/logr"
	"strconv"
)

type TableInfo struct {
	TableName string
	TableType interface{}
	TimeFlag  int8 //1:一个时间缀  2:两个时间缀
	OprFlag   int8 //1:增量   2:全量
	CrtName   string
	UpdName   string
	TimName   string //一个时间缀
	NeedFormat bool  //日期是否需要格式化
}

type Tbl_table_sync struct {
	TABLE_NAME string
	CUR_CRT_TS string
	CUR_UPD_TS string
	CUR_TIME   string
}

func (t Tbl_table_sync) TableName() string {
	return "TBL_TABLE_SYNC"
}

type Tbl_leaguer struct {
	LEAGUER_NO     string `gorm:"primary_key"`
	LEAGUER_NAME   string
	LEAGUER_TYPE   string
	LEAGUER_INFO   string
	LEAGUER_STATUS int64
	REC_CRT_TS     time.Time
	REC_UPD_TS     time.Time
}

func (t Tbl_leaguer) TableName() string {
	return "TBL_LEAGUER"
}

func (t *Tbl_leaguer) BeforeCreate(tx *gorm.DB) error {

	switch t.LEAGUER_TYPE {
	case "BCEFG":
		t.LEAGUER_TYPE = "BL"
	case "BCEFGHI":
		t.LEAGUER_TYPE = "BCDF"
	}
	return nil
}

func (t *Tbl_leaguer) BeforeUpdate(tx *gorm.DB) error {

	switch t.LEAGUER_TYPE {
	case "BCEFG":
		t.LEAGUER_TYPE = "BL"
	case "BCEFGHI":
		t.LEAGUER_TYPE = "BCDF"
	}
	return nil
}

////////////////////////////////////
type Tbl_user struct {
	USER_ID              int64 `gorm:"AUTO_INCREMENT"`
	LEAGUER_NO           string
	USER_NAME            string `gorm:"primary_key"`
	AUTH_KEY             string
	PASSWORD_HASH        string
	PASSWORD_RESET_TOKEN string
	EMAIL                string
	USER_TYPE            string
	USER_INFO            string
	USER_STATUS          int64
	USER_NOTICE          string
	REC_CRT_TS           time.Time
	REC_UPD_TS           time.Time
}

func (t Tbl_user) TableName() string {
	return "TBL_USER"
}

func (t *Tbl_user) BeforeCreate() error  {
	t.USER_ID = 0
	return nil
}

func (t *Tbl_user) AfterCreate(tx *gorm.DB) error {
	logr.Info("insert userid: ", t.USER_ID)
	authAss := Tbl_auth_assignment{}
	authAss.User_id = strconv.FormatInt(t.USER_ID, 10)
	authAss.Item_name = "商户查询员"
	authAss.Created_at = time.Now().Unix()
	return tx.Create(authAss).Error
}

////////////////////////////////////
type Tbl_auth_assignment struct {
	Item_name  string `xorm:"ITEM_NAME" gorm:"primary_key"`
	User_id    string `xorm:"USER_ID"  gorm:"primary_key"`
	Created_at int64  `xorm:"CREATED_AT"`
}

func (t Tbl_auth_assignment) TableName() string {
	return "TBL_AUTH_ASSIGNMENT"
}

////////////////////////////////////
type Tbl_auth_item struct {
	Name        string `xorm:"NAME"  gorm:"primary_key"`
	Type        int64  `xorm:"TYPE"`
	Description string `xorm:"DESCRIPTION"`
	Rule_name   string `xorm:"RULE_NAME"`
	Data        string `xorm:"ITEM_DATA"`
	Created_at  int64  `xorm:"CREATED_AT"`
	Updated_at  int64  `xorm:"UPDATED_AT"`
}

func (t Tbl_auth_item) TableName() string {
	return "TBL_AUTH_ITEM"
}

////////////////////////////////////
type Tbl_auth_item_child struct {
	Parent string `xorm:"PARENT" gorm:"primary_key"`
	Child  string `xorm:"CHILD" gorm:"primary_key"`
}

func (t Tbl_auth_item_child) TableName() string {
	return "TBL_AUTH_ITEM_CHILD"
}

////////////////////////////////////
type Tbl_menu struct {
	Id     int64  `xorm:"ID"  gorm:"primary_key"`
	Name   string `xorm:"NAME"`
	Parent int64  `xorm:"PARENT"`
	Route  string `xorm:"MENU_ROUTE"`
	Order  int64  `xorm:"MENU_ORDER"`
	Data   string `xorm:"MENU_DATA"`
}

func (t Tbl_menu) TableName() string {
	return "TBL_MENU"
}

////////////////////////////////////
type Tbl_user_info struct {
	USER_TYPE string `gorm:"primary_key"`
	USER_DESC string
	USER_INFO string
}

func (t Tbl_user_info) TableName() string {
	return "TBL_USER_INFO"
}

///////////////////////////////////
type Tbl_clear_txn struct {
	COMPANY_CD          string
	INS_ID_CD           string
	ACQ_INS_ID_CD       string
	FWD_INS_ID_CD       string
	MCHT_CD             string
	MCHT_NAME           string
	MCHT_SHORT_NAME     string
	MCC_CD              string
	MCC_CD_42           string
	MCC_DESC            string
	TRANS_DATE_TIME     string
	STLM_DATE           string
	TRANS_KIND          string
	TXN_DESC            string
	TRANS_STATE         string
	STLM_FLG            string
	TRANS_AMT           string
	CREDITCARDLIMIT     string
	CUP_SSN             string
	AUTHR_ID_RESP       string
	PAN                 string
	CARD_KIND_DIS       string
	BANK_CODE           string
	BANK_NAME           string
	BRANCH_CD           string
	BRANCH_NM           string
	TERM_ID             string
	ORG_TRANS_DATE_TIME string
	ORG_CUP_SSN         string
	POS_ENTRY_MODE      string
	RSP_CODE            string
	TRUE_FEE_MOD        string
	TRUE_FEE_BI         string
	TRUE_FEE_FD         string
	TRUE_FEE_FFD        string
	VAR_1               string
	VAR_2               string
	VAR_3               string
	VAR_4               string
	VIR_FEE_MOD         string
	VIR_FEE_BI          string
	VIR_FEE_BD          string
	VIR_FEE_FD          string
	MCHT_FEE            string
	VAR_5               string
	MCHT_VIR_FEE        string
	STAND_BANK_FEE      string
	BANK_FEE            string
	HZJG_FEE            string
	JGSY                string
	AIP_FEE             string `gorm:"column:AIP_FEE"`
	MCHT_SET_AMT        string
	HZJGYFPPFWF         string
	JGYFPPFWF           string
	AIPYFPPFWF          string `gorm:"column:AIPYFPPFWF"`
	ERR_FEE_IN          string
	ERR_FEE_OUT         string
	ERR_CODE            string
	JT_MCHT_CD          string
	EXPAND_ORG_CD       string
	SPE_SERV_INST       string
	PROP_INS            string
	EXPAND_ORG_FEE      string
	SPE_SERV_FEE        string
	PROP_INS_FEE        string
	EXPAND_ORG_PP       string
	SPE_SERV_PP         string
	PROP_INS_PP         string
	EXPAND_FEE_IN       string
	EXPAND_FEE_OUT      string
	CUP_IFINSIDE_SIGN   string `gorm:"column:CUP_IFINSIDE_SIGN"`
	SP_CHARG_TYPE       string
	SP_CHARG_LEV        string
	TERM_SSN            string
	SN_SSN              string
	UP_CHL_ID           string
	CONV_MCHT_CD        string
	CONV_TERM_ID        string
	CHL_TRUE_FEE        string
	CHL_STD_FEE         string
	CHL_FEE_PRE_FLG     string
	SYS_SER             string
	VAR_6               string
	QUDAO_FEE           string
	QUDAO_FEE_MIN       string
	QUDAO_FEE_MIX       string
	QUDAO_FEE_FD        string
	INS_FEE             string
	INS_MY_FEE          string
	INS_COST_FEE        string
	INS_MY_FEE_AMT      string
	INS_SPLIT_FEE       string
	INS_RES_FEE         string
	PINP_FEE            string
	PINP_FEE_INF        string
	PINP_FEE_TOP        string
	PINP_STAT           string
	T0_STAT             string
	KEY_RSP             string `gorm:"primary_key"`
	REMARK              string
	REMARK1             string `gorm:"column:REMARK1"`
	REMARK2             string `gorm:"column:REMARK2"`
	REMARK3             string `gorm:"column:REMARK3"`
	REMARK4             string `gorm:"column:REMARK4"`
	REMARK5             string `gorm:"column:REMARK5"`
}

func (t Tbl_clear_txn) TableName() string {
	return "TBL_CLEAR_TXN"
}

func (t *Tbl_clear_txn) BeforeCreate() error {
	t.KEY_RSP = t.KEY_RSP + KeyRspFix
	return nil
}

////////////////////////////////////////
type Tbl_dictionaryitem struct {
	DIC_TYPE	string `gorm:"primary_key"`
	DIC_CODE	string `gorm:"primary_key"`
	DIC_NAME	string
	DISP_ORDER	string
	MEMO	string
	UPDATE_TIME	time.Time
}

func (t Tbl_dictionaryitem) TableName() string {
	return "TBL_DICTIONARYITEM"
}

////////////////////////////////////////
type Tbl_ins_ctrl_inf struct {
	INS_ID_CD	string `gorm:"primary_key"`
	INS_COMPANY_CD	string
	PROD_CD	string `gorm:"primary_key"`
	BIZ_CD	string `gorm:"primary_key"`
	CTRL_STA	string
	INS_BEG_TM	string
	INS_END_TM	string
	MSG_RESV_FLD1	string `gorm:"column:MSG_RESV_FLD1"`
	MSG_RESV_FLD2	string `gorm:"column:MSG_RESV_FLD2"`
	MSG_RESV_FLD3	string `gorm:"column:MSG_RESV_FLD3"`
	MSG_RESV_FLD4	string `gorm:"column:MSG_RESV_FLD4"`
	MSG_RESV_FLD5	string `gorm:"column:MSG_RESV_FLD5"`
	MSG_RESV_FLD6	string `gorm:"column:MSG_RESV_FLD6"`
	MSG_RESV_FLD7	string `gorm:"column:MSG_RESV_FLD7"`
	MSG_RESV_FLD8	string `gorm:"column:MSG_RESV_FLD8"`
	MSG_RESV_FLD9	string `gorm:"column:MSG_RESV_FLD9"`
	MSG_RESV_FLD10	string `gorm:"column:MSG_RESV_FLD10"`
	REC_OPR_ID	string
	REC_UPD_OPR	string
	REC_CRT_TS	time.Time
	REC_UPD_TS	time.Time
}
func (t Tbl_ins_ctrl_inf) TableName() string {
	return "TBL_INS_CTRL_INF"
}

////////////////////////////////////////
type Tbl_ins_inf struct {
	INS_ID_CD	string `gorm:"primary_key"`
	INS_COMPANY_CD	string
	INS_TYPE	string
	INS_NAME	string
	INS_PROV_CD	string
	INS_CITY_CD	string
	INS_REGION_CD	string
	INS_STA	string
	INS_STLM_TP	string
	INS_ALO_STLM_CYCLE	string
	INS_ALO_STLM_MD	string
	INS_STLM_C_NM	string
	INS_STLM_C_ACCT	string
	INS_STLM_C_BK_NO	string
	INS_STLM_C_BK_NM	string
	INS_STLM_D_NM	string
	INS_STLM_D_ACCT	string
	INS_STLM_D_BK_NO	string
	INS_STLM_D_BK_NM	string
	MSG_RESV_FLD1	string `gorm:"column:MSG_RESV_FLD1"`
	MSG_RESV_FLD2	string `gorm:"column:MSG_RESV_FLD2"`
	MSG_RESV_FLD3	string `gorm:"column:MSG_RESV_FLD3"`
	MSG_RESV_FLD4	string `gorm:"column:MSG_RESV_FLD4"`
	MSG_RESV_FLD5	string `gorm:"column:MSG_RESV_FLD5"`
	MSG_RESV_FLD6	string `gorm:"column:MSG_RESV_FLD6"`
	MSG_RESV_FLD7	string `gorm:"column:MSG_RESV_FLD7"`
	MSG_RESV_FLD8	string `gorm:"column:MSG_RESV_FLD8"`
	MSG_RESV_FLD9	string `gorm:"column:MSG_RESV_FLD9"`
	MSG_RESV_FLD10	string `gorm:"column:MSG_RESV_FLD10"`
	REC_OPR_ID	string
	REC_UPD_OPR	string
	REC_CRT_TS	time.Time
	REC_UPD_TS	time.Time
}
func (t Tbl_ins_inf) TableName() string {
	return "TBL_INS_INF"
}

////////////////////////////////////////
type Tbl_mcht_biz_deal struct {
	MCHT_CD	string `gorm:"primary_key"`
	PROD_CD	string `gorm:"primary_key"`
	BIZ_CD	string `gorm:"primary_key"`
	TRANS_CD	string `gorm:"primary_key"`
	OPER_IN	string
	REC_OPR_ID	string
	REC_UPD_OPR	string
	REC_CRT_TS	time.Time
	REC_UPD_TS	time.Time
}
func (t Tbl_mcht_biz_deal) TableName() string {
	return "TBL_MCHT_BIZ_DEAL"
}

////////////////////////////////////////
type Tbl_mcht_inf struct {
	MCHT_CD	string `gorm:"primary_key"`
	SN	string
	AIP_BRAN_CD	string `gorm:"column:AIP_BRAN_CD"`
	GROUP_CD	string
	ORI_CHNL	string
	ORI_CHNL_DESC	string
	BANK_BELONG_CD	string
	DVP_BY	string
	MCC_CD_18	string
	APPL_DATE	string
	UP_BC_CD	string
	UP_AC_CD	string
	UP_MCC_CD	string
	NAME	string
	NAME_BUSI	string
	BUSI_LICE_NO	string
	BUSI_RANG	string
	BUSI_MAIN	string
	CERTIF	string
	CERTIF_TYPE	string
	CERTIF_NO	string
	NATION_CD	string
	PROV_CD	string
	CITY_CD	string
	AREA_CD	string
	REG_ADDR	string
	CONTACT_NAME	string
	CONTACT_PHONENO	string
	ISGROUP	string
	MONEYTOGROUP	string
	STLM_WAY	string
	STLM_WAY_DESC	string
	STLM_INS_CIRCLE	string
	APPR_DATE	time.Time
	STATUS	string
	DELETE_DATE	time.Time
	UC_BC_CD_32	string
	K2WORKFLOWID	string `gorm:"column:K2WORKFLOWID"`
	SYSTEMFLAG	string
	APPROVALUSERNAME	string
	FINALARRPOVALUSERNAME	string
	IS_UP_STANDARD	string
	BILLINGTYPE	string
	BILLINGLEVEL	string
	SLOGAN	string
	EXT1	string `gorm:"column:EXT1"`
	EXT2	string `gorm:"column:EXT2"`
	EXT3	string `gorm:"column:EXT3"`
	EXT4	string `gorm:"column:EXT4"`
	AREA_STANDARD	string
	MCHTCD_AREA_CD	string
	UC_BC_CD_AREA	string
	REC_OPR_ID	string
	REC_UPD_OPR	string
	REC_CRT_TS	time.Time
	REC_UPD_TS	time.Time
	OPER_IN	string
	REC_APLLY_TS	time.Time
	OEM_ORG_CODE	string
}
func (t Tbl_mcht_inf) TableName() string {
	return "TBL_MCHT_INF"
}

////////////////////////////////////////
type Tbl_prod_biz_trans_map struct {
	PROD_CD	string `gorm:"primary_key"`
	BIZ_CD	string `gorm:"primary_key"`
	TRANS_CD	string `gorm:"primary_key"`
	UPDATE_DATE	string
	DESCRIPTION	string `gorm:"column:DESCRIPTION"`
	RESV_FLD1	string `gorm:"column:RESV_FLD1"`
	RESV_FLD2	string `gorm:"column:RESV_FLD2"`
	RESV_FLD3	string `gorm:"column:RESV_FLD3"`
}
func (t Tbl_prod_biz_trans_map) TableName() string {
	return "TBL_PROD_BIZ_TRANS_MAP"
}

////////////////////////////////////////
type Tbl_term_inf struct {
	MCHT_CD	string `gorm:"primary_key"`
	TERM_ID	string `gorm:"primary_key"`
	TERM_TP	string
	BELONG	string
	BELONG_SUB	string
	TMNL_MONEY_INTYPE	string
	TMNL_MONEY	int64
	TMNL_BRAND	string
	TMNL_MODEL_NO	string
	TMNL_BARCODE	string
	DEVICE_CD	string
	INSTALLLOCATION	string
	TMNL_INTYPE	string
	DIAL_OUT	string
	DEAL_TYPES	string
	REC_OPR_ID	string
	REC_UPD_OPR	string
	REC_CRT_TS	time.Time
	REC_UPD_TS	time.Time
	APP_CD	string
	SYSTEMFLAG	string
	STATUS	string
	ACTIVE_CODE	string
}
func (t Tbl_term_inf) TableName() string {
	return "TBL_TERM_INF"
}

////////////////////////////////////////
type Tbl_tfr_his_trn_log struct {
	TRANS_DT string
	TRANS_MT string
	SRC_QID int64 `gorm:"column:SRC_QID"`
	DES_QID int64 `gorm:"column:DES_QID"`
	MA_TRANS_CD string
	MA_TRANS_NM string
	KEY_RSP string `gorm:"primary_key"`
	KEY_REVSAL string
	KEY_CANCEL string
	RESP_CD string
	TRANS_ST string
	MA_TRANS_SEQ int64
	ORIG_MA_TRANS_SEQ int64
	ORIG_TRANS_SEQ string
	ORIG_TERM_SEQ string
	ORIG_TRANS_DT string
	MA_SETTLE_DT string `gorm:"column:MA_SETTLE_DT"`
	ACCESS_MD string
	MSG_TP string
	PRI_ACCT_NO string
	ACCT_TP string
	TRANS_PROC_CD string
	TRANS_AT string
	TRANS_TD_TM string
	TERM_SEQ string
	ACPT_TRANS_TM string
	ACPT_TRANS_DT string
	MCHNT_TP string
	POS_ENTRY_MD_CD string
	POS_COND_CD string
	ACPT_INS_ID_CD string
	FWD_INS_ID_CD string
	TERM_ID string
	MCHNT_CD string
	CARD_ACCPTR_NM string
	RETRI_REF_NO string
	REQ_AUTH_ID string
	TRANS_SUBCATA string
	INDUSTRY_ADDN_INF string
	TRANS_CURR_CD string
	SEC_CTRL_INF string
	IC_DATA string
	UDF_FLD_PURE string
	CERTIF_ID string
	NETWORK_MGMT_INF_CD string
	ORIG_DATA_ELEMNT string
	RCV_INS_ID_CD string
	TFR_IN_ACCT_NO_PURE string
	TFR_IN_ACCT_TP string
	TFR_OUT_ACCT_NO_PURE string
	ACPT_INS_RESV_PURE string
	TRR_OUT_ACCT_TP string
	ISS_INS_ID_CD string
	CARD_ATTR string
	CARD_CLASS string
	CARD_MEDIA string
	CARD_BIN string
	CARD_BRAND string
	ROUT_INS_ID_CD string
	ACPT_REGION_CD string
	BUSS_REGION_CD string
	USR_NO_TP string
	USR_NO_REGION_CD string
	USR_NO_REGION_ADDN_CD string
	USR_NO string
	SP_INS_ID_CD string
	INDUSTRY_INS_ID_CD string
	ROUT_INDUSTRY_INS_ID_CD string
	INDUSTRY_MCHNT_CD string
	INDUSTRY_TERM_CD string
	INDUSTRY_MCHNT_TP string
	ENTRUST_TP string
	PMT_MD string
	PMT_TP string
	PMT_NO string
	PMT_MCHNT_CD string
	PMT_NO_INDUSTRY_INS_ID_CD string
	PRI_ACCT_NO_CONV string
	TRANS_AT_CONV string
	TRANS_DT_TM_CONV string
	TRANS_SEQ_CONV string
	MCHNT_TP_CONV string
	RETRI_REF_NO_CONV string
	ACPT_INS_ID_CD_CONV string
	TERM_ID_CONV string
	MCHNT_CD_CONV string
	MCHNT_NM_CONV string
	UDF_FLD_PURE_CONV string
	SP_INS_ID_CD_CONV string
	EXPIRE_DT string
	SETTLE_DT string `gorm:"column:SETTLE_DT"`
	TRANS_FEE string
	RESP_AUTH_ID string
	ACPT_RESP_CD string
	ADDN_RESP_DATA_PURE string
	ADDN_AT_PURE string
	ISS_ADDN_DATA_PURE string
	IC_RES_DAT_CUPS string
	SW_RESV_PURE string
	ISS_INS_RESV_PURE string
	INDUSTRY_RESP_CD string
	DEBT_AT string
	DTL_INQ_DATA string
	TRANS_CHNL string
	INTERCH_MD_CD string
	TRANS_CHK_IN string
	MCHT_STLM_FLG string
	INS_STLM_FLG string
	MSG_RESV_FLD1 string `gorm:"column:MSG_RESV_FLD1"`
	MSG_RESV_FLD2 string `gorm:"column:MSG_RESV_FLD2"`
	MSG_RESV_FLD3 string `gorm:"column:MSG_RESV_FLD3"`
	TRANS_MTH int64
	REC_UPD_TS time.Time
	REC_CRT_TS time.Time
	PROD_CD string
	TRAN_TP string
	BIZ_CD string
	REVEL_FLG string
	CANCEL_FLG string
	MSG_RESV_FLD4 string `gorm:"column:MSG_RESV_FLD4"`
	MSG_RESV_FLD5 string `gorm:"column:MSG_RESV_FLD5"`
	MSG_RESV_FLD6 string `gorm:"column:MSG_RESV_FLD6"`
	MSG_RESV_FLD7 string `gorm:"column:MSG_RESV_FLD7"`
	MSG_RESV_FLD8 string `gorm:"column:MSG_RESV_FLD8"`
	MSG_RESV_FLD9 string `gorm:"column:MSG_RESV_FLD9"`
}
func (t Tbl_tfr_his_trn_log) TableName() string {
	return "TBL_TFR_HIS_TRN_LOG"
}

func (t *Tbl_tfr_his_trn_log) BeforeCreate() error {
	t.KEY_RSP = t.KEY_RSP + KeyRspFix
	return nil
}

////////////////////////////////////////
type Tbl_tfr_pre_auth_log Tbl_tfr_his_trn_log

func (t Tbl_tfr_pre_auth_log) TableName() string {
	return "TBL_TFR_PRE_AUTH_LOG"
}

func Clear(v interface{})  {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
