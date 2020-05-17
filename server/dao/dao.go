package dao

import (
	redigo "github.com/gomodule/redigo/redis"
	"mk-api/server/conf"
)

// 全局Redis连接电池
var Redis = new(RedisPools)

// 全局mysql db
var Db = new(DB)

// Register your redis-cli pool here
type RedisPools struct {
	tokenPool *redigo.Pool
}

func init() {
	Db = NewMySQL(&conf.C.MysqlWrite, &conf.C.MysqlRead)
	Redis.tokenPool = NewRedisPool(&conf.C.RedisToken)
}
