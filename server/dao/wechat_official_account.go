package dao

import (
	"strconv"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"mk-api/server/conf"
)

func NewOfficialAccount() *officialaccount.OfficialAccount {
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        conf.C.RedisWechat.Host + ":" + strconv.Itoa(conf.C.RedisWechat.Port),
		Password:    conf.C.RedisWechat.Password,
		Database:    conf.C.RedisWechat.Db,
		MaxIdle:     conf.C.RedisWechat.MaxIdle,
		MaxActive:   conf.C.RedisWechat.MaxActive,
		IdleTimeout: conf.C.RedisWechat.IdleTimeout,
	}
	redisCache := cache.NewRedis(redisOpts)
	cfg := &offConfig.Config{
		AppID:          conf.C.WeChat.AppID,
		AppSecret:      conf.C.WeChat.AppSecret,
		Token:          conf.C.WeChat.Token,
		EncodingAESKey: conf.C.WeChat.EncodingAESKey,
		Cache:          redisCache,
	}
	return wc.GetOfficialAccount(cfg)
}
