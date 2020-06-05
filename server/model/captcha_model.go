package model

import (
	"github.com/gomodule/redigo/redis"
	"mk-api/server/dao"
)

type CaptchaModel interface {
	Save(key string, val string) (err error)
	Check(key string, val2Check string) (ok bool)
}

type captchaDatabase struct {
	redisPool *redis.Pool
}

func (db *captchaDatabase) Check(key string, val2Check string) (ok bool) {
	cli := db.redisPool.Get()
	defer cli.Close()
	valStored, err := redis.String(cli.Do("GET", key))
	if err != nil || valStored != val2Check {
		return false
	}
	return true
}

func (db *captchaDatabase) Save(key string, val string) (err error) {
	cli := db.redisPool.Get()
	defer cli.Close()
	_, err = cli.Do("SETEX", key, 180, val)
	return
}

// model 层有错误要抛出去给 service 层
func NewCaptchaModel() CaptchaModel {
	return &captchaDatabase{
		redisPool: dao.Rdb.TokenRdbP,
	}
}
