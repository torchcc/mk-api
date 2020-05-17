package dao

import (
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"mk-api/server/conf"
)

func NewRedisPool(conf *conf.RedisConfig) *redigo.Pool {
	server := conf.Host + ":" + strconv.Itoa(conf.Port)
	maxIdl, maxActive := 1, 5
	if conf.MaxIdle != 0 {
		maxIdl = conf.MaxIdle
	}
	if conf.MaxActive != 0 {
		maxActive = conf.MaxActive
	}

	p := &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
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
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     maxIdl,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	}

	return p
}
