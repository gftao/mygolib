package goxorm

import (
	"fmt"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	//_ "github.com/mattn/go-oci8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/widuu/goini"
	"mygolib/modules/config"
	"os"
	"strconv"
	"time"
 )

type Enginesql struct {
	Engin   *xorm.Engine
	Session *xorm.Session
	Rows    *xorm.Rows
}

type levelmap map[int]map[string]map[string]string

var err error

func opendb(file string) (*Enginesql, error) {
	//检测配置文件
	var Engine Enginesql
	talkChan := loadConfig(file)
	dbdata := make(map[string]string)
	dbdata["dbtype"] = talkChan[0]["db"]["dbtype"]
	dbdata["dbname"] = talkChan[0]["db"]["dbname"]
	dbdata["dbuser"] = talkChan[0]["db"]["dbuser"]
	dbdata["dbpass"] = talkChan[0]["db"]["dbpass"]
	dbdata["dbpost"] = talkChan[0]["db"]["dbpost"]
	dbdata["dbip"] = talkChan[0]["db"]["dbip"]
	dbdata["dblog"] = talkChan[0]["db"]["dblog"]
	dbdata["idlcon"] = talkChan[0]["db"]["idlcon"]
	dbdata["maxcon"] = talkChan[0]["db"]["maxcon"]
	var a *xorm.Engine
	var dbstr string
	switch dbdata["dbtype"] {
	case "mysql":
		//dbstr = "hypmpcld" + ":" + "hypmpcld" + "@tcp(" + "192.168.20.77" + ":" + "3306" + ")/" + "hypmpcld" +
		//	"?charset=utf8&parseTime=True&loc=Local"
		dbstr = dbdata["dbuser"] + ":" + dbdata["dbpass"] + "@tcp(" + dbdata["dbip"] + ":" + dbdata["dbpost"] + ")/" + dbdata["dbname"] +
			"?charset=utf8&parseTime=True&loc=Local"
		fmt.Println(dbstr)
		a, err := xorm.NewEngine("mysql", dbstr)
		if err != nil {
			return &Engine, fmt.Errorf("连接失败 %s,%s\n", err, dbstr)
		}
		Engine.Engin = a
	default:
		//gopmp/gopmp@192.168.20.12/tfrbke
		dbstr = dbdata["dbuser"] + "/" + dbdata["dbpass"] + "@" + dbdata["dbip"] + "/" + dbdata["dbname"]
		fmt.Println(dbstr)
		a, err = xorm.NewEngine("oci8", dbstr)
		if err != nil {
			return &Engine, fmt.Errorf("连接失败 %s,%s\n", err, dbstr)
		}
		Engine.Engin = a
	}
	Engine.Engin.SetTableMapper(core.SameMapper{})
	Engine.Engin.SetColumnMapper(core.SameMapper{})
	idlCon, err := strconv.Atoi(dbdata["idlcon"])
	if err != nil || idlCon <= 0 {
		fmt.Println("空闲连接默认10")
		idlCon = 10
	}

	maxCon, err := strconv.Atoi(dbdata["maxcon"])
	if err != nil || maxCon <= 0 {
		fmt.Println("最大连接默认100")
		maxCon = 100
	}

	Engine.Engin.SetMaxIdleConns(idlCon)
	Engine.Engin.SetMaxOpenConns(maxCon)
	//Engine.Engin.SetColumnMapper(core.GonicMapper{})

	var ok bool
	ok, err = strconv.ParseBool(dbdata["dblog"])
	if err != nil {
		fmt.Println("数据库日志模式取失败，默认不打印", err)
		ok = false
	}
	if !ok {
		nilLogger := xorm.DiscardLogger{}
		Engine.Engin.SetLogger(nilLogger)
	}
	Engine.Engin.ShowSQL(ok)

	return &Engine, nil
}

func loadConfig(file string) (talkChan levelmap) {
	talkChan = make(levelmap)
	conf := goini.SetConfig(file)
	talkChan1 := conf.ReadList()

	for k, v := range talkChan1 {
		talkChan[k] = v
		for k1, v1 := range v {
			talkChan[k][k1] = v1
			for k2, v2 := range v1 {
				talkChan[k][k1][k2] = v2
			}

		}
	}

	return talkChan
}

/*创建Session*/
func (m *Enginesql) CreateSession() (*xorm.Session, error) {
	sess := m.Engin.NewSession()
	err := sess.Begin()
	if err != nil {
		sess.Close()
	}
	return sess, err
}

/*提交Session*/
func (m *Enginesql) CommitSession(sess *xorm.Session) error {
	err := sess.Commit()
	if err != nil {
		sess.Rollback()
	}
	sess.Close()
	return err
}

/*回滚Session*/
func (m *Enginesql) RollbackSession(sess *xorm.Session) error {
	err := sess.Rollback()
	sess.Close()
	return err
}

func (m *Enginesql) Query(sql string) ([]map[string][]uint8, error) {
	results, err := m.Engin.Query(sql)
	return results, err
}
func (m *Enginesql) Exec(sql string, args ...interface{}) (interface{}, error) {
	results, err := m.Engin.Exec(sql, args...)
	return results, err
}
func (m *Enginesql) Begin() error {
	m.Session = m.Engin.NewSession()
	err := m.Session.Begin()
	return err
}
func (m *Enginesql) End() {
	m.Session.Close()
}
func (m *Enginesql) Rollback() {
	m.Session.Rollback()
}

func (m *Enginesql) FindOne(links interface{}, querystring string, args ...interface{}) (bool, error) {
	has, err := m.Engin.Where(querystring, args...).Get(links)
	return has, err
}
func (m *Enginesql) FindSql(links interface{}, querystring string, args ...interface{}) error {
	err := m.Engin.Sql(querystring, args...).Find(links)
	return err
}
func (m *Enginesql) FindAll(links interface{}, querystring string, args ...interface{}) error {
	err := m.Engin.Where(querystring, args...).Find(links)
	return err
}
func (m *Enginesql) NewFetch(links interface{}, querystring string, args ...interface{}) error {
	var err error
	m.Rows, err = m.Engin.Where(querystring, args...).Rows(links)
	return err
}
func (m *Enginesql) FetchEnd() {
	m.Rows.Close()
}
func (m *Enginesql) FetchRow() bool {
	r := m.Rows.Next()
	return r
}
func (m *Enginesql) FetchToStruct(links interface{}) error {
	err := m.Rows.Scan(links)
	return err
}
func (m *Enginesql) StructToInsert(links interface{}) error {
	_, err = m.Session.Insert(links)
	return err
}
func (m *Enginesql) StructToUpdate(links interface{}, querystring string, args ...interface{}) error {
	_, err = m.Session.Where(querystring, args...).Update(links)
	return err
}
func (m *Enginesql) StructToDelete(links interface{}, querystring string, args ...interface{}) error {
	_, err = m.Session.Where(querystring, args...).Delete(links)
	return err
}
func (m *Enginesql) Commit() error {
	err := m.Session.Commit()
	return err
}
func (m *Enginesql) SqlToExec(sql string, args ...interface{}) error {
	_, err = m.Session.Exec(sql, args...)
	return err
}
func (m *Enginesql) Close() error {
	return m.Engin.Close()
}
func (m *Enginesql) GetsysDate() (t time.Time, e error) {
	var sysdate []byte
	//results, err := m.Engin.Query("select systimestamp from dual")
	results, err := m.Engin.Query("select now()")
	if err != nil {
		return t, err
	}
	for _, v := range results {
		for _, sysdate = range v {
		}
	}
	s, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", byteString(sysdate))

	return s, err
}

func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func (m *Enginesql) LogWriter(file string) error {
	m.Engin.ShowSQL(true)
	m.Engin.Logger().SetLevel(core.LOG_DEBUG)
	//var f *os.File
	//if _, err := os.Stat(file); os.IsNotExist(err) {
	//	//f, err = os.OpenFile(file, os.O_CREATE)
	//	f, err = os.Create(file)
	//} else {
	//	f, err = os.OpenFile(file, os.O_APPEND, 0666)
	//}
	//if err != nil {
	//	fmt.Printf("打开文件出错%s\n", err)
	//}
	f, err := os.OpenFile(file,
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		os.ModePerm|os.ModeTemporary)

	m.Engin.SetLogger(xorm.NewSimpleLogger(f))
	return err
}
func GetInstance() *Enginesql {
	err := instance.Engin.Ping()
	dbconf := config.StringDefault("dbconf", "")

	for err != nil {
		fmt.Errorf("数据库连接已经断开，重新连接 %s", err)
		instance, err = opendb(dbconf)
		time.Sleep(time.Duration(5) * time.Second)
	}
	return instance
}
