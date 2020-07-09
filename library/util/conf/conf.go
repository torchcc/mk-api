package conf

import (
	"mk-api/library/superconf"
)

var C *Config = nil

type Config struct {
	Cos                    cosConfig
	QiniuCos               qiniuConfig
	RegisterSmsMsgTemplate smsMsgTemplateConfig
	Local                  superconf.Config
}

type cosConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
}

type smsMsgTemplateConfig struct {
	SmsSdkAppid string
	TemplateID  string // zk 读取
}

type qiniuConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	ImgPath   string `json:"img_path"`
}

// first define your conf data structure above here , second register your configs here
func init() {
	cfg := Config{}

	var allConfigs = make(map[string]interface{})
	allConfigs["/superconf/third_party/cos"] = &cfg.Cos
	allConfigs["/superconf/third_party/sms/register_msg_template"] = &cfg.RegisterSmsMsgTemplate
	allConfigs["/superconf/third_party/qiniu"] = &cfg.QiniuCos

	sc := superconf.NewSuperConfig(&allConfigs)
	cfg.Local = *(sc.Config)
	C = &cfg
}
