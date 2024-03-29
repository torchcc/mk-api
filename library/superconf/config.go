package superconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"mk-api/deployment"
)

type SuperConfig struct {
	Config   *Config
	zkConn   *zk.Conn
	Path2Cfg *map[string]interface{}
}

type Config struct {
	Deploy    string   `json:"deploy"`
	ZKServers []string `json:"zk_servers"`
}

func NewSuperConfig(configs *map[string]interface{}) (s *SuperConfig) {

	cfg, err := loadJsonFile()
	if err != nil {
		panic("load" + deployment.SUPERCONF_JSON_ABS_PATH + " failed: " + err.Error())
	}

	var validCfg = Config{
		Deploy:    cfg.Env.Deploy,
		ZKServers: strings.Split(cfg.Env.Zookeeper.Host, ","),
	}

	conn, _, err := zk.Connect(validCfg.ZKServers, time.Second)
	if err != nil {
		panic("failed to connect to zk, err is " + err.Error())
	}

	for k, v := range *configs {
		watch(conn, k, v)
	}

	return &SuperConfig{
		Config:   &validCfg,
		zkConn:   conn,
		Path2Cfg: configs,
	}
}

// structure of superconf.json
type superconfJson struct {
	Env struct {
		Deploy    string `json:"deploy"`
		Zookeeper struct {
			Host string `json:"host"`
		} `json:"zookeeper"`
	} `json:"env"`
}

func loadJsonFile() (cfg *superconfJson, err error) {
	data, err := ioutil.ReadFile(deployment.SUPERCONF_JSON_ABS_PATH)
	if err != nil {
		return
	}
	cfg = &superconfJson{}
	err = json.Unmarshal(data, &cfg)
	return
}

func (sc *SuperConfig) connect() (conn *zk.Conn) {
	var err error
	if conn, _, err = zk.Connect(sc.Config.ZKServers, time.Second); err != nil {
		panic(err)
	}
	return
}

// TODO re-init mysql db and redis-cli pool when their confs change
// TODO map is not concurrency safe, use the map in sync/atomic instead
func watch(conn *zk.Conn, path string, key interface{}) {

	dates := make(chan []byte)
	errCh := make(chan error)

	go func() {
		for {
			t, _, events, err := conn.GetW(path)
			if err != nil {
				errCh <- err
				return
			}
			dates <- t
			evt := <-events
			if evt.Err != nil {
				errCh <- evt.Err
				return
			}
			fmt.Println("GetW ")
		}
	}()

	go func() {
		for {
			select {
			case data := <-dates:
				switch key.(type) {
				case *int:
					t, _ := key.(*int)
					*t, _ = strconv.Atoi(string(data[:]))
				case *string:
					t, _ := key.(*string)
					*t = string(data[:])
				default:
					_ = json.Unmarshal(data, &key)
				}
			case err := <-errCh:
				fmt.Printf("watchStr error %+v\n", err)
				conn.Close()
				return
			}
		}
	}()

	data, _, _ := conn.Get(path)
	switch key.(type) {
	case *int:
		t, _ := key.(*int)
		*t, _ = strconv.Atoi(string(data[:]))
	case *string:
		t, _ := key.(*string)
		*t = string(data[:])
	default:
		_ = json.Unmarshal(data, &key)
	}

}
