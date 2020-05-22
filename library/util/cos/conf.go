package cos

import (
	"mk-api/library/superconf"
)

var C *Config = nil

type Config struct {
	Cos   cosConfig
	Local superconf.Config
}

type cosConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

// first define your conf data structure above here , second register your configs here
func init() {
	cfg := Config{}

	var allConfigs = make(map[string]interface{})
	allConfigs["/superconf/third_party/cos"] = &cfg.Cos

	sc := superconf.NewSuperConfig(&allConfigs)
	cfg.Local = *(sc.Config)
	C = &cfg
}
