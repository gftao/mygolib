package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	"golib/modules/config"
	"golib/modules/goxorm"
	"golib/modules/logr"
	"golib/tools/tableSync/models"
	"os"
	"time"
)

var dbc *goxorm.Enginesql

func main() {

	err := config.InitModuleByParams("./export.ini")

	if err != nil {
		fmt.Println("初始化配置文件失败", err)
		return
	}

	err = goxorm.InitModel()
	if err != nil {
		fmt.Println("数据库初始化失败", err)
		return
	}

	err = logr.InitModules()
	if err != nil {
		fmt.Println("日志初始化失败", err)
		return
	}

	if len(os.Args) >1 {
		models.ExecTime = os.Args[1]
	}
	if len(os.Args) > 2 {
		models.FilePrefix = os.Args[2]
	}

	//NLS_LANG=American_america.zhs16gbk
	//NLS_TIMESTAMP_TZ_FORMAT=YYYY-MM-DD HH24:MI:SS.FF6 TZR
	//NLS_TIMESTAMP_FORMAT=YYYY-MM-DD HH24:MI:SS.FF6
	//NLS_DATE_FORMAT=YYYYMMDD
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	os.Setenv("NLS_TIMESTAMP_TZ_FORMAT", "YYYY-MM-DD HH24:MI:SS.FF6 TZR")
	os.Setenv("NLS_TIMESTAMP_FORMAT", "YYYY-MM-DD HH24:MI:SS.FF6")
	os.Setenv("NLS_DATE_FORMAT", "YYYYMMDD")

	//err = gormdb.InitModule()
	//if err != nil {
	//	fmt.Println("数据库初始化失败", err)
	//	return
	//}

	dbc = goxorm.GetInstance()
	dbc.Engin.ShowSQL(true)
	tx, err := dbc.CreateSession()
	if err != nil {
		logr.Error("创建session失败", err)
		return
	}
	defer tx.Close()
	err = expTables(tx)
	if err != nil {
		logr.Error("导出表失败", err)
		tx.Rollback()
		return
	}
	tx.Commit()

	logr.Info("导出完成")
}

func expTables(tx *xorm.Session) error {

	for _, t := range models.TableListExp {
		var begTime, maxTime string
		var begCrtTIme, begUpdTime, maxCrtTime, maxUpdTime string
		var allFlag bool
		tableSync := models.Tbl_table_sync{}

		if t.TimeFlag == models.ONETIME {

			resultes, err := tx.Query("select max(\"" + t.TimName + "\") as T from \"" + t.TableName + "\"")
			if err != nil {
				logr.Error("查询最大时间失败", err)
				return err
			}
			l, _ := resultes[0]["T"]
			logr.Debug("get time:" + string(l))
			if t.NeedFormat {
				tp, err := time.Parse(models.TIMEINFORMAT, string(l))
				if err != nil {
					logr.Error("格式化时间失败", string(l), err)
					return err
				}
				maxTime = tp.Format(models.TIMEOUTFORMAT)
			} else {
				maxTime = string(l)
			}
		} else if t.TimeFlag == models.TWOTIME {
			resultes, err := tx.Query("select max(\"" + t.CrtName + "\") as T from \"" + t.TableName + "\"")
			if err != nil {
				logr.Error("查询最大时间失败", err)
				return err
			}
			l, _ := resultes[0]["T"]
			mt, err := time.Parse(models.TIMEINFORMAT, string(l))
			maxCrtTime = mt.Format(models.TIMEOUTFORMAT)
			resultes, err = tx.Query("select max(\"" + t.UpdName + "\") as T from \"" + t.TableName + "\"")
			if err != nil {
				logr.Error("查询最大时间失败", err)
				return err
			}
			l, _ = resultes[0]["T"]
			mt, err = time.Parse(models.TIMEINFORMAT, string(l))
			maxUpdTime = mt.Format(models.TIMEOUTFORMAT)
		}
		ok, err := dbc.FindOne(&tableSync, " table_name = ? ", t.TableName)
		if !ok && err == nil {
			logr.Info("同步记录未找到，全量同步" + t.TableName)
			allFlag = true
			if t.TimeFlag == models.ONETIME {
				tableSync.CUR_TIME = maxTime
				tableSync.TABLE_NAME = t.TableName
				i, err := tx.InsertOne(&tableSync)
				if i != 1 || err != nil {
					logr.Error("插入记录失败", i, err)
					return err
				}
			} else if t.TimeFlag == models.TWOTIME {
				tableSync.CUR_CRT_TS = maxCrtTime
				tableSync.CUR_UPD_TS = maxUpdTime
				tableSync.TABLE_NAME = t.TableName
				i, err := tx.InsertOne(&tableSync)
				if i != 1 || err != nil {
					logr.Error("插入记录失败", i, err)
					return err
				}
			}
		} else {
			allFlag = false
			if t.TimeFlag == models.ONETIME {
				begTime = tableSync.CUR_TIME
				tableSync.CUR_TIME = maxTime
				if begTime != maxTime {
					_, err = tx.Where(" table_name = ? ", t.TableName).
						Update(models.Tbl_table_sync{CUR_TIME:tableSync.CUR_TIME})
				}
			} else if t.TimeFlag == models.TWOTIME {
				begCrtTIme = tableSync.CUR_CRT_TS
				begUpdTime = tableSync.CUR_UPD_TS
				tableSync.CUR_CRT_TS = maxCrtTime
				tableSync.CUR_UPD_TS = maxUpdTime
				if begCrtTIme != maxCrtTime || begUpdTime != maxUpdTime {
					_, err = tx.Where(" table_name = ? ", t.TableName).Update(models.Tbl_table_sync{
						CUR_CRT_TS:tableSync.CUR_CRT_TS, CUR_UPD_TS:tableSync.CUR_UPD_TS})
				}
			}
			if err != nil {
				logr.Error("更新同步日期失败", err)
				return err
			}
		}

		//开始数据导出
		logr.Info("vars: ", begTime, maxTime, begCrtTIme, begUpdTime, maxCrtTime, maxUpdTime, allFlag)

		fpath := models.FilePrefix + models.ExecTime + "/"
		os.MkdirAll(fpath, 0777)
		fname := fpath + t.TableName + "_" + models.ExecTime + models.FileEdx
		fp, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logr.Error("打开文件失败"+fname, err)
			return err
		}

		//从配置文件取where条件
		selWhere := config.StringDefault(t.TableName, "")
		if selWhere != "" {
			tx.Where(selWhere)
			logr.Info("添加条件控制" + selWhere)
		}

		var rows *xorm.Rows
		if t.TimeFlag == models.ONETIME {
			if allFlag {
				rows, err = tx.Where("\""+t.TimName+"\" <= ? or \""+t.TimName+"\" is null ",
					maxTime).Rows(t.TableType)
			} else {
				rows, err = tx.Where("\""+t.TimName+"\" <= ? and \""+t.TimName+"\" >= ? ",
					maxTime, begTime).Rows(t.TableType)
			}
		} else if t.TimeFlag == models.TWOTIME {
			if allFlag {
				rows, err = tx.Where("\""+t.CrtName+"\" <= ? or \""+t.UpdName+"\" <= ?  or \""+
					t.CrtName+"\" is null or \""+t.UpdName+"\" is null ",
					maxCrtTime, maxUpdTime).Rows(t.TableType)
			} else {
				rows, err = tx.Where("(\""+t.CrtName+"\" <= ? and \""+t.CrtName+"\" >= ?) or ( \""+
					t.UpdName+"\" <= ? and \""+t.UpdName+"\" >= ? ) ",
					maxCrtTime, begCrtTIme, maxUpdTime, begUpdTime).Rows(t.TableType)
			}
		} else if t.TimeFlag == models.NOTIME {
			rows, err = tx.Rows(t.TableType)
		}
		if err != nil {
			logr.Error("数据库操作失败", err)
			fp.Close()
			return err
		}
		defer rows.Close()
		lineNum := 0
		for rows.Next() {
			models.Clear(t.TableType)
			err = rows.Scan(t.TableType)
			if err != nil {
				logr.Error("取记录失败", err)
				fp.Close()
				return err
			}
			by, err := json.Marshal(t.TableType)
			if err != nil {
				logr.Error("转换json失败", err)
				return err
			}
			fp.Write(by)
			fp.WriteString("\n")
			lineNum++
		}
		logr.Infof("表[%s]操作完成，共[%d]条", t.TableName, lineNum)

		fp.Close()

	}
	return nil
}
