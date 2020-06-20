package conf

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// this also can be the test for superconf
func TestConf(t *testing.T) {

	go func() {
		for {
			b, err := json.Marshal(C)
			if err != nil {
				t.Errorf("json.Marshal(%v) error(%v)", C, err)
			}
			fmt.Printf("configures as follow: %v\n", string(b[:]))
			time.Sleep(time.Second)
		}
	}()

	// for {
	// 	xtime.Sleep(xtime.Second)
	// }
	time.Sleep(time.Second * 20)
}
