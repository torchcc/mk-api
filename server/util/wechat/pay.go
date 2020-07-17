package wechat

import (
	"mk-api/server/conf"
	"strconv"
	"strings"
	"time"

	"github.com/silenceper/wechat/v2/pay/order"
	"github.com/silenceper/wechat/v2/util"
)

func Launch2ndPay(nonceStr string, prepayId string) (cfg *order.Config, err error) {
	cfg = &order.Config{}
	var (
		buffer    strings.Builder
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	)
	const signType = "MD5"

	buffer.WriteString("appId=")
	buffer.WriteString(conf.C.WeChat.AppID)
	buffer.WriteString("&nonceStr=")
	buffer.WriteString(nonceStr)
	buffer.WriteString("&package=")
	buffer.WriteString("prepay_id=" + prepayId)
	buffer.WriteString("&signType=")
	buffer.WriteString(signType)
	buffer.WriteString("&timeStamp=")
	buffer.WriteString(timestamp)
	buffer.WriteString("&key=")
	buffer.WriteString(conf.C.WeChat.PayKey)
	sign, err := util.CalculateSign(buffer.String(), signType, conf.C.WeChat.PayKey)
	if err != nil {
		return
	}

	// 签名
	cfg.PaySign = sign
	cfg.NonceStr = nonceStr
	cfg.Timestamp = timestamp
	cfg.PrePayID = prepayId
	cfg.SignType = signType
	cfg.Package = "prepay_id=" + prepayId
	return
}
