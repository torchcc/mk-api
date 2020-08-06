package dao

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"mk-api/server/conf"
)

// Redis redis cache
type Redis struct {
	conn *redis.Pool
}

func NewRedis(conf *conf.RedisConfig) *Redis {
	server := conf.Host + ":" + strconv.Itoa(conf.Port)
	maxIdl, maxActive := 1, 5
	if conf.MaxIdle != 0 {
		maxIdl = conf.MaxIdle
	}
	if conf.MaxActive != 0 {
		maxActive = conf.MaxActive
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			c.Do("SELECT", conf.Db)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     maxIdl,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	}

	return &Redis{conn: pool}
}

func NewRedisPool(conf *conf.RedisConfig) *redis.Pool {
	fmt.Printf("conf is %#v, creating redis pool.....", conf)
	server := conf.Host + ":" + strconv.Itoa(conf.Port)
	maxIdl, maxActive := 1, 5
	if conf.MaxIdle != 0 {
		maxIdl = conf.MaxIdle
	}
	if conf.MaxActive != 0 {
		maxActive = conf.MaxActive
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					fmt.Printf("NewRedisPool failed, params is [%v], err is [%s]", conf, err.Error())
					panic("failed to create redis pool !")
					return nil, err
				}
			}
			c.Do("SELECT", conf.Db)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				fmt.Printf("NewRedisPool failed, params is [%v], err is [%s]", conf, err.Error())
				panic("failed to PING redis pool !")
			}
			return err
		},
		MaxIdle:     maxIdl,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	}

	return pool
}

// SetConn 设置conn
func (r *Redis) SetConn(conn *redis.Pool) {
	r.conn = conn
}

// Get 获取一个值
func (r *Redis) Get(key string) ([]byte, error) {
	conn := r.conn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// SetEx 设置一个值 a
func (r *Redis) Set(key string, val interface{}) (err error) {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}

	_, err = conn.Do("SET", key, data)

	return
}

// SetEx 设置一个值 a
func (r *Redis) SetEx(key string, val interface{}, timeout time.Duration) (err error) {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}

	_, err = conn.Do("SETEX", key, int64(timeout/time.Second), data)

	return
}

// IsExist 判断key是否存在
func (r *Redis) Exists(key string) bool {
	conn := r.conn.Get()
	defer conn.Close()

	a, _ := conn.Do("EXISTS", key)
	i := a.(int64)
	if i > 0 {
		return true
	}
	return false
}

// Delete 删除
func (r *Redis) Delete(key string) error {
	conn := r.conn.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}

	return nil
}

// SetEx 设置一个值 a
func (r *Redis) HSet(key string, field string, val interface{}) (err error) {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}

	_, err = conn.Do("HSET", key, field, data)

	return
}

func (r *Redis) HGet(key string, field string) interface{} {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	var err error
	if data, err = redis.Bytes(conn.Do("HGET", key, field)); err != nil {
		return nil
	}
	var reply interface{}
	if err = json.Unmarshal(data, &reply); err != nil {
		return nil
	}

	return reply
}

func (r *Redis) Expire(key string, timeout time.Duration) (err error) {
	conn := r.conn.Get()
	defer conn.Close()

	_, err = conn.Do("EXPIRE", key, int64(timeout/time.Second))

	return
}

// LikeDeletes batch delete
func (r *Redis) LikeDeletes(key string) error {
	conn := r.conn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = r.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
