package server

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	Log.WithFields(logrus.
		Fields{"order_id": 123456, "user_id": 1}).
		Errorf("订单付款失败: err: %s", "服务器错误")
}
