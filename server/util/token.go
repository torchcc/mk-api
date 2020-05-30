package util

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func OpenId2Token(openId string) string {
	h := md5.New()
	strconv.FormatInt(time.Now().Unix(), 10)
	h.Write([]byte(openId + strconv.FormatInt(time.Now().Unix(), 10)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
