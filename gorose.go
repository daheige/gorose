package gorose

import (
	"database/sql"
	"errors"
	"github.com/gohouse/gorose/drivers"
)

var (
	// Connect is the Connection Object
	Connect Connection
)

// Connection is the database pre handle
type Connection struct {
	// DB is origin DB
	DB *sql.DB
	// all config sets
	DbConfig map[string]interface{}
	// default database
	Default string
	// current config on use
	CurrentConfig map[string]string
	//// all sql logs
	//SqlLog []string
	//// if in transaction, the code auto change
	//Trans bool
	// max open connections
	SetMaxOpenConns int
	// max freedom connections leave
	SetMaxIdleConns int
}

func init() {
	Connect.SetMaxOpenConns = 0
	Connect.SetMaxIdleConns = -1
}

// Open instance of sql.DB.Oper
// if args has 1 param , it will be derect connection or with default config set
// if args has 2 params , the second param will be the default dirver key
func Open(args ...interface{}) (Connection, error) {
	//fmt.Println(args)
	//return Connect, errors.New("dsf")
	if len(args) == 1 {
		// continue
	} else if len(args) == 2 {
		if confReal, ok := args[1].(string); ok {
			Connect.Default = confReal
		} else {
			// 指定默认数据库只能为字符串!
			return Connect, errors.New("only str allowed of default database name")
		}
	} else {
		// Open方法只接收1个或2个参数!
		return Connect, errors.New("1 or 2 params need in Open() method")
	}
	// 解析config
	err := Connect.parseConfig(args[0])
	if err != nil {
		return Connect, err
	}

	// 驱动数据库
	errs := Connect.boot()

	return Connect, errs
}

// Parse input config
func (conn *Connection) parseConfig(args interface{}) error {
	if confReal, ok := args.(map[string]string); ok { // direct connection
		Connect.CurrentConfig = confReal
	} else if confReal, ok := args.(map[string]interface{}); ok {
		// store the full connection
		Connect.DbConfig = confReal
		// if set the Default conf, store it
		if defaultDb, ok := confReal["Default"]; ok {
			// judge if seted
			if Connect.Default == "" {
				Connect.Default = defaultDb.(string)
			}
		}
		if Connect.Default == "" {
			// 配置文件默认数据库链接未设置
			return errors.New("the default database is missing in config!")
		}
		// 获取指定的默认数据库链接信息
		var connections map[string]map[string]string
		if connectionsInterface, ok := confReal["Connections"]; ok {
			switch connectionsInterface.(type) {
			case map[string]map[string]string:
				connections = connectionsInterface.(map[string]map[string]string)
			case map[string]interface{}:
				connectionsTmp := connectionsInterface.(map[string]interface{})
				if connectionsTmpReal, ok := connectionsTmp[Connect.Default]; ok {
					switch connectionsTmpReal.(type) {
					case map[string]string:
						connections = map[string]map[string]string{Connect.Default: connectionsTmpReal.(map[string]string)}
					default:
						return errors.New("the database connections format error !")
					}
				}
			default:
				return errors.New("the database connections format error !")
			}
		} else {
			return errors.New("the database connections missing !")
		}
		if defaultDbConnection, ok := connections[Connect.Default]; ok {
			Connect.CurrentConfig = defaultDbConnection
		} else {
			// 指定的数据库链接不存在!
			return errors.New("the database for using is missing!")
		}
		// 设置连接池信息
		if mo, ok := confReal["SetMaxOpenConns"]; ok {
			if moInt, ok := mo.(int); ok {
				conn.SetMaxOpenConns = moInt
			} else {
				// 连接池信息配置的值只能是数字
				return errors.New("the value of connection pool config need int")
			}
		}
		if mi, ok := confReal["SetMaxIdleConns"]; ok {
			if miInt, ok := mi.(int); ok {
				conn.SetMaxIdleConns = miInt
			} else {
				return errors.New("the value of connection pool config need int")
			}
		}
	} else {
		return errors.New("format error in database config!")
	}
	return nil
}

// Boot sql driver
func (conn *Connection) boot() error {
	//dbObj := Connect.CurrentConfig
	var driver, dsn string
	var err error

	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	driver,dsn = drivers.GetDsnByDriverName(Connect.CurrentConfig)

	// 开始驱动
	conn.DB, err = sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	conn.DB.SetMaxOpenConns(conn.SetMaxOpenConns)
	conn.DB.SetMaxIdleConns(conn.SetMaxIdleConns)

	// 检查是否可以ping通
	err2 := conn.DB.Ping()

	return err2
}

// Close database
func (conn *Connection) Close() error {
	//Connect.SqlLog = []string{}
	return conn.DB.Close()
}

// Ping db
func (conn *Connection) Ping() error {
	return conn.DB.Ping()
}

// Table is set table from database
func (conn *Connection) Table(table string) *Database {
	return conn.GetInstance().Table(table)
}

//// Begin transaction begin
//func (conn *Connection) Begin() {
//	Tx, _ = DB.Begin()
//	Connect.Trans = true
//}
//
//// Commit is transaction commit
//func (conn *Connection) Commit() {
//	Tx.Commit()
//	Connect.Trans = false
//}
//
//// Rollback is transaction rollback
//func (conn *Connection) Rollback() {
//	Tx.Rollback()
//	Connect.Trans = false
//}
//
//// Transaction is simple transaction
//func (conn *Connection) Transaction(closure func() error) bool {
//	//defer func() {
//	//	if err := recover(); err != nil {
//	//		conn.Rollback()
//	//		panic(err)
//	//	}
//	//}()
//
//	conn.Begin()
//	err := closure()
//	if err != nil {
//		conn.Rollback()
//		return false
//	}
//	conn.Commit()
//
//	return true
//}

// Query str
func (conn *Connection) Query(args ...interface{}) ([]map[string]interface{}, error) {
	return conn.GetInstance().Query(args...)
}

// Execute str
func (conn *Connection) Execute(args ...interface{}) (int64, error) {
	return conn.GetInstance().Execute(args...)
}

// GetInstance , get the database object
func (conn *Connection) GetInstance() *Database {
	//var database *Database
	//return database
	return &Database{}
}

// JsonEncode : parse json
func (conn *Connection) JsonEncode(arg interface{}) string {
	return conn.GetInstance().JsonEncode(arg)
}

//// LastSql is get last query sql
//func (conn *Connection) LastSql() string {
//	if len(Connect.SqlLog) > 0 {
//		return Connect.SqlLog[len(Connect.SqlLog)-1:][0]
//	}
//	return ""
//}
//
//// SqlLogs is all sql query logs in this request
//func (conn *Connection) SqlLogs() []string {
//	return Connect.SqlLog
//}

// GetDB is get origin *sql.DB
func (conn *Connection) GetDB() *sql.DB {
	return conn.DB
}
