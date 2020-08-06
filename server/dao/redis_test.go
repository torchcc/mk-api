package dao

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"mk-api/server/conf"
)

var redisCfg = &conf.RedisConfig{
	Host:        "106.53.124.190",
	Port:        6979,
	Db:          0,
	Timeout:     0,
	MaxIdle:     0,
	MaxActive:   0,
	IdleTimeout: 0,
	Password:    "",
}

func TestRedis(t *testing.T) {

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

func TestEncapsulatedRedis(t *testing.T) {
	rd := NewRedis(redisCfg)

	if err := rd.Set("f", "xue mei"); err != nil {
		t.Errorf("redis set err: %v", err)
	} else {
		t.Logf("redis set done!")
	}

	err := rd.SetEx("name", "lisa", time.Second*10)
	if err != nil {
		t.Errorf("redis err: %v", err)
	} else {
		t.Logf("SetEx Done!")
	}

	if err := rd.HSet("xuemei", "lvl", 20); err != nil {
		t.Errorf("redis Hset err: %v", err)
	} else {
		t.Logf("HSet Done!")
	}

	res := rd.HGet("xuemei", "lvl")
	t.Logf("Hget Done! the val of xuemei cup is %v", res.(float64))
}
