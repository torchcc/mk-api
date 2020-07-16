package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"mk-api/server/conf"
)

// 全局Redis连接电池
var Rdb = new(Rdbs)

// 全局mysql db
var Db *sqlx.DB

// 全局go-cache
var GoCache *cache.Cache

// 全局 OfficialAccount
var AffAcc *officialaccount.OfficialAccount

// Register your redis-cli
type Rdbs struct {
	TokenRdb  *Redis
	TokenRdbP *redis.Pool
}

func init() {
	Db = NewMySQLx(&conf.C.MysqlWrite)
	Rdb.TokenRdb = NewRedis(&conf.C.RedisToken)
	Rdb.TokenRdbP = NewRedisPool(&conf.C.RedisToken)
	GoCache = NewGoCache()
	AffAcc = NewOfficialAccount()
}
