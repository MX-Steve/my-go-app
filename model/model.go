// 数据库初始化，建立连接
package model

import (
	"database/sql"
	"fmt"

	"github.com/MX-Steve/my-go-app/vip"

	_ "github.com/go-sql-driver/mysql"
)

type DbWorker struct {
	dsn string
	Db  *sql.DB
}

func NewDb() DbWorker {
	user := vip.GetIniData("mysql.user")
	port := vip.GetIniData("mysql.port")
	host := vip.GetIniData("mysql.host")
	password := vip.GetIniData("mysql.password")
	db := vip.GetIniData("mysql.db")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", user, password, host, port, db)
	dbw := DbWorker{dsn: dsn}
	dbtemp, _ := sql.Open("mysql", dbw.dsn)
	dbw.Db = dbtemp
	return dbw
}
