package conf

import (
	"fmt"

	"mk-api/library/superconf"
)

const ServiceName = "mk-server"

var C *Config

type Config struct {
	MysqlRead   MysqlConfig
	MysqlWrite  MysqlConfig
	RedisToken  RedisConfig
	RedisWechat RedisConfig
	Local       superconf.Config
	MongoLog    MongoConfig
	WeChat      WechatConfig
	// GenerateOrderKafka kafka.Config
	RecvOpenIds []string // 运营人员open列表
}

type MysqlConfig struct {
	KeepConnectionAlive bool   `json:"keepConnectionAlive"`
	MaxConnections      int    `json:"maxConnections"`
	Port                int    `json:"port"`
	MinFreeConnections  int    `json:"minFreeConnections"`
	Host                string `json:"host"`
	Database            string `json:"database"`
	Password            string `json:"password"`
	User                string `json:"user"`
	Charset             string `json:"charset"`
}

type RedisConfig struct {
	Port        int    `json:"port"`
	Db          int    `json:"db"`
	Timeout     int    `json:"timeout"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
	Password    string `json:"password"`
	Host        string `json:"host"`
}

type MongoConfig struct {
	Port             int    `json:"port"`
	Host             string `json:"host"`
	Db               string `json:"db"`
	AuthDb           string `json:"auth_db"`
	ConnectionString string `json:"connection_string"`
	Collection       string `json:"collection"`
	Password         string `json:"password"`
	User             string `json:"user"`
}

type WechatConfig struct {
	AppID          string `json:"app_id"`
	AppSecret      string `json:"app_secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	PayMchID       string `json:"pay_mch_id"`     // 支付 - 商户 Id
	PayNotifyURL   string `json:"pay_notify_url"` // 支付 - 接受微信支付结果通知的接口地址
	PayKey         string `json:"pay_key"`        // 支付 - 商户后台设置的支付 key
}

// first define your conf data structure above here , second register your configs here
func init() {
	cfg := Config{}

	var allConfigs = make(map[string]interface{})
	allConfigs["/superconf/union/mysql/read"] = &cfg.MysqlRead
	allConfigs["/superconf/union/mysql/write"] = &cfg.MysqlWrite
	allConfigs["/superconf/union/redis/token"] = &cfg.RedisToken
	allConfigs["/superconf/union/redis/wechat"] = &cfg.RedisWechat
	allConfigs["/superconf/union/mongo/log"] = &cfg.MongoLog
	allConfigs["/superconf/third_party/wechat"] = &cfg.WeChat
	allConfigs["/superconf/third_party/receiver_open_ids"] = &cfg.RecvOpenIds

	sc := superconf.NewSuperConfig(&allConfigs)
	cfg.Local = *(sc.Config)
	C = &cfg
	fmt.Printf("all config in server is : %#v", cfg)
}
