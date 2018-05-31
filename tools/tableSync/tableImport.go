package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"golib/modules/config"
	"golib/modules/gormdb"
	"golib/modules/logr"
	"golib/tools/tableSync/models"
	"os"
	"strings"
)

func main() {

	err := config.InitModuleByParams("./import.ini")

	if err != nil {
		fmt.Println("初始化配置文件失败", err)
		return
	}

	err = logr.InitModules()
	if err != nil {
		fmt.Println("日志初始化失败", err)
		return
	}


	models.KeyRspFix = config.StringDefault("KeyRspFix", "")
	logr.Info("取出后缀KeyRspFix:" + models.KeyRspFix)

	err = gormdb.InitModule()
	if err != nil {
		fmt.Println("数据库初始化失败", err)
		return
	}

	if len(os.Args) >1 {
		models.ExecTime = os.Args[1]
	}
	if len(os.Args) > 2 {
		models.FilePrefix = os.Args[2]
	}

	tx := gormdb.GetInstance().Begin()
	if err != nil {
		logr.Error("创建session失败", err)
		return
	}
	defer tx.Close()
	err = impTable(tx)
	if err != nil {
		logr.Error("导入表失败", err)
		tx.Rollback()
		return
	}
	tx.Commit()

	logr.Info("导入完成")
}

func impTable(tx *gorm.DB) error {

	for _, t := range models.TableListImp {
		logr.Info("开始处理" + t.TableName)

		fname := models.FilePrefix + models.ExecTime + "/" + t.TableName + "_" + models.ExecTime + models.FileEdx
		fp, err := os.Open(fname)
		if err != nil {
			logr.Error("打开文件失败" + fname)
			return err
		}

		//if t.OprFlag == models.ALLTYPE {
		//	err = tx.Delete(t.TableType).Error
		//	if err != nil {
		//		logr.Warn("全量同步删除数据失败", err)
		//	}
		//}

		scan := bufio.NewScanner(fp)
		lineNum := 0
		for scan.Scan() {
			line := scan.Bytes()
			models.Clear(t.TableType)
			err = json.Unmarshal(line, t.TableType)
			if err != nil {
				logr.Errorf("解析json失败, 文件[%s], 行数[%d], %v", scan, line, err)
				fp.Close()
				return err
			}

			err = tx.Create(t.TableType).Error
			if err != nil {
				if strings.Contains(err.Error(), "Duplicate") {
					//tx.Delete(t.TableType)
					err = tx.Model(t.TableType).Update(t.TableType).Error
					if err != nil {
						logr.Error("已经存在数据更新失败", t.TableType, err)
						return err
					}
				} else {
					logr.Error("添加记录失败", err)
					return err
				}
			}
			lineNum++
		}
		logr.Infof("文件[%s]处理完成[%d]", fname, lineNum)
		fp.Close()
	}

	return nil
}

