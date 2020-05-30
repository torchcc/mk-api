package dao

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"mk-api/server/conf"
	"mk-api/server/util"
)

func NewMySQLx(c *conf.MysqlConfig) *sqlx.DB {
	dsn := c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.Database + "?charset=utf8mb4&autocommit=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		util.Log.Panicf("sqlx 数据库驱动连接不到mysql服务器: %v\n", err.Error())
	}

	db.SetMaxOpenConns(c.MaxConnections)
	db.SetMaxIdleConns(c.MinFreeConnections)
	return db
}
