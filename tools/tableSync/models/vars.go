package models

import "time"

var TableListExp = []TableInfo{
	//{"TBL_AUTH_ASSIGNMENT", &Tbl_auth_assignment{}, ONETIME, ADDTYPE, "", "", "CREATED_AT", false},
	{"TBL_AUTH_ITEM", &Tbl_auth_item{}, ONETIME, ADDTYPE, "", "", "CREATED_AT", false},
	{"TBL_AUTH_ITEM_CHILD", &Tbl_auth_item_child{}, NOTIME, ALLTYPE, "", "", "", false},
	/////////////{"TBL_menu", &Tbl_menu{}, NOTIME, ALLTYPE, "", "", ""},
	//////////
	{"TBL_LEAGUER", &Tbl_leaguer{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_CLEAR_TXN", &Tbl_clear_txn{}, ONETIME, ADDTYPE, "", "", "STLM_DATE", false},
	{"TBL_DICTIONARYITEM", &Tbl_dictionaryitem{}, ONETIME, ADDTYPE, "", "", "UPDATE_TIME", true},
	{"TBL_INS_CTRL_INF", &Tbl_ins_ctrl_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_INS_INF", &Tbl_ins_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_MCHT_BIZ_DEAL", &Tbl_mcht_biz_deal{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_MCHT_INF", &Tbl_mcht_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_PROD_BIZ_TRANS_MAP", &Tbl_prod_biz_trans_map{}, ONETIME, ADDTYPE, "", "", "UPDATE_DATE", false},
	{"TBL_TERM_INF", &Tbl_term_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_TFR_HIS_TRN_LOG", &Tbl_tfr_his_trn_log{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_TFR_PRE_AUTH_LOG", &Tbl_tfr_pre_auth_log{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_USER", &Tbl_user{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	///////////////{"TBL_USER_INFO", &Tbl_user_info{}, NOTIME, ADDTYPE, "", "", "", false},
}

var TableListImp = []TableInfo{
	//{"TBL_AUTH_ASSIGNMENT", &Tbl_auth_assignment{}, ONETIME, ADDTYPE, "", "", "CREATED_AT", false},
	{"TBL_AUTH_ITEM", &Tbl_auth_item{}, ONETIME, ADDTYPE, "", "", "created_at", false},
	//{"TBL_auth_item_child", &Tbl_auth_item_child{}, NOTIME, ALLTYPE, "", "", "", false},
	/////////////{"TBL_menu", &Tbl_menu{}, NOTIME, ALLTYPE, "", "", ""},
	//////////
	{"TBL_LEAGUER", &Tbl_leaguer{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_CLEAR_TXN", &Tbl_clear_txn{}, ONETIME, ADDTYPE, "", "", "STLM_DATE", false},
	{"TBL_DICTIONARYITEM", &Tbl_dictionaryitem{}, ONETIME, ADDTYPE, "", "", "UPDATE_TIME", true},
	{"TBL_INS_CTRL_INF", &Tbl_ins_ctrl_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_INS_INF", &Tbl_ins_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_MCHT_BIZ_DEAL", &Tbl_mcht_biz_deal{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_MCHT_INF", &Tbl_mcht_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_PROD_BIZ_TRANS_MAP", &Tbl_prod_biz_trans_map{}, ONETIME, ADDTYPE, "", "", "UPDATE_DATE", false},
	{"TBL_TERM_INF", &Tbl_term_inf{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_TFR_HIS_TRN_LOG", &Tbl_tfr_his_trn_log{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_TFR_PRE_AUTH_LOG", &Tbl_tfr_pre_auth_log{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	{"TBL_USER", &Tbl_user{}, TWOTIME, ADDTYPE, "REC_CRT_TS", "REC_UPD_TS", "", true},
	///////////////{"TBL_USER_INFO", &Tbl_user_info{}, NOTIME, ADDTYPE, "", "", "", false},
}


const TIMEINFORMAT = "2006-01-02T15:04:05.999999+08:00"
const TIMEOUTFORMAT = "2006-01-02 15:04:05.999999"
var ExecTime = time.Now().Format("20060102")

var KeyRspFix = ""

var FilePrefix = "./files/"
var FileEdx = ".dat"

const (
	ONETIME = 1
	TWOTIME = 2
	NOTIME  = 3

	ADDTYPE = 1
	ALLTYPE = 2
)
