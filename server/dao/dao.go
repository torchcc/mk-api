package dao

import (
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"mk-api/server/conf"
)

// 全局Redis连接电池
var Redis = new(RedisPools)

// 全局mysql db
var Db *sqlx.DB

// Register your redis-cli pool here
type RedisPools struct {
	tokenPool *redigo.Pool
}

func init() {
	Db = NewMySQLx(&conf.C.MysqlWrite)
	Redis.tokenPool = NewRedisPool(&conf.C.RedisToken)
}
