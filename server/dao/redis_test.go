package dao

import (
	"testing"

	"github.com/gomodule/redigo/redis"
	"mk-api/server/conf"
)

func TestRedis(t *testing.T) {

	redisCfg := &conf.RedisConfig{
		Host:        "106.53.124.190",
		Port:        6979,
		Db:          0,
		Timeout:     0,
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: 0,
		Password:    "",
	}
	// if you want to test whether the global configuration conf.C , use the row below
	// redisCfg := &conf.C.RedisToken
	redisP := NewRedisPool(redisCfg)
	c := redisP.Get()
	defer c.Close()

	_, err := c.Do("Set", "troy", 123)
	if err != nil {
		t.Errorf("redis pool E	rr: %v\n", err)
		return
	}

	r, err := redis.Int(c.Do("Get", "troy"))
	if err != nil {
		t.Errorf("redis pool Err: %v\n", err)
		return
	}
	t.Logf("key: %s, value: %d\n", "troy", r)

}
