package gormdb

import (
	"github.com/jinzhu/gorm"
	"mygolib/modules/config"
	"mygolib/modules/myLogger"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"fmt"
)

var instance *gorm.DB
var initFlg bool = false

func InitModule() error {

	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}

	var err error
	instance, err = initDb()

	if err != nil {
		return err
	}

	initFlg = true

	return nil
}

func HasModuleInit() bool {
	return initFlg
}

func initDb() (*gorm.DB, error) {

	var err error
	var db *gorm.DB

	var connStr string
	config.SetSection("database")
	dbtype := config.StringDefault("type", "mysql")
	dbhost := config.StringDefault("host", "127.0.0.1")
	dbport := config.StringDefault("port", "3306")
	dbname := config.StringDefault("dbname", "prodPmpCld")
	dbuser := config.StringDefault("user", "root")
	dbpasswd := config.StringDefault("passwd", "")

	logModel := config.BoolDefault("logmodel", false)
	idlCon := config.IntDefault("idlcon", 10)
	maxCon := config.IntDefault("maxcon", 100)
	sslModel := config.StringDefault("sslModel", "disable")

	switch dbtype {
	case "mysql":
		//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		connStr = dbuser + ":" + dbpasswd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname +
			"?charset=utf8&parseTime=True&loc=Local"
		db, err = gorm.Open("mysql", connStr)
		myLogger.Info("conn to db:" + connStr)
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	case "postgres":
		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			dbhost, dbport, dbuser, dbname, sslModel, dbpasswd)
		myLogger.Info("conn to post db:" + connStr)
		db, err = gorm.Open("postgres", connStr)
	default:
		myLogger.Error("非法的数据库类型" + dbtype)
		return nil, errors.New("非法的数据库类型" + dbtype)
	}

	if err != nil {
		return nil, err
	}

	db.LogMode(logModel)
	db.DB().Ping()
	db.DB().SetMaxIdleConns(idlCon)
	db.DB().SetMaxOpenConns(maxCon)

	//db.SetLogger(myLogger )

	return db, nil
}

func GetInstance() *gorm.DB {
	err := instance.DB().Ping()
	for err != nil {
		myLogger.Error("数据库连接已经断开，重新连接", err)
		instance, err = initDb()
		time.Sleep(time.Duration(5) * time.Second)
	}
	return instance
}
