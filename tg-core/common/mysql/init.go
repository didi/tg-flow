package mysql

import (
	"database/sql"
	"log"
	"time"

	"github.com/didi/tg-flow/tg-core/conf"

	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

var Handler *sql.DB

const (
	SECTION = "mysql"
)

func InitMySqlHandler() {
	log.Println("tg-core mysql init start...")
	sec, err := conf.Handler.GetSection(SECTION)
	if err != nil {
		log.Fatal("MySQL init error: ", err)
	}
	host := sec.GetStringMust("host", "localhost")
	port := sec.GetIntMust("port", 3306)
	user := sec.GetStringMust("user", "root")
	pwd := sec.GetStringMust("password", "")
	charset := sec.GetStringMust("charset", "utf8")
	dbName := sec.GetStringMust("db_name", "")
	timeout := sec.GetIntMust("conn_timeout", 500)
	readTimeout := sec.GetIntMust("read_timeout", 500)
	writeTimeout := sec.GetIntMust("write_timeout", 500)
	connMaxLifetime := sec.GetIntMust("conn_max_lifetime", 0)
	maxIdleConns := sec.GetIntMust("max_idle_conns", 2)
	maxOpenConns := sec.GetIntMust("max_open_conns", 2)

	Handler, err = manager.New(dbName, user, pwd, host).
		Set(
			manager.SetCharset(charset),
			manager.SetTimeout(time.Duration(timeout)*time.Millisecond),
			manager.SetReadTimeout(time.Duration(readTimeout)*time.Millisecond),
			manager.SetWriteTimeout(time.Duration(writeTimeout)*time.Millisecond),
			manager.SetParseTime(true),
			manager.SetLoc("Local"),
		).Port(int(port)).Open(true)
	if err != nil {
		log.Fatal("Init mysql error: ", err)
	}

	Handler.SetConnMaxLifetime(time.Duration(connMaxLifetime))
	Handler.SetMaxIdleConns(int(maxIdleConns))
	Handler.SetMaxOpenConns(int(maxOpenConns))
	log.Println("tg-core mysql init successful!")
}
