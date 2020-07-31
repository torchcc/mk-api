package wechat

import (
	"testing"
	"time"

	"mk-api/server/conf"
)

func TestPush(t *testing.T) {
	// 测试付款成功后发送给员工
	// testOrderPaidNotifyStaff(t)

	// 付款成功后发送给客户
	// testOrderPaidNotifyClient(t)

	// 发起退款后通知员工
	// testRefundLaunchedNotifyStaff(t)

	// 测试预约成功后发送给客户
	// testAppointmentMadeNotifyClient(t)
}

func testOrderPaidNotifyStaff(t *testing.T) {
	openIds := []string{"oDvnPw4zKAmraE2eccSUHinSya5E"}
	var amount float64 = 56.89
	OrderPaidNotifyStaff(openIds, "12345676", amount, time.Now().Unix())
}

func testOrderPaidNotifyClient(t *testing.T) {

	openId := "oDvnPw4zKAmraE2eccSUHinSya5E"
	outTradeNo := "2112465451521"
	// orderTime := "2020年07月28日 19:21:21"
	var amount float64 = 99.99
	var orderId int64 = 25
	OrderPaidNotifyClient(openId, outTradeNo, amount, orderId, time.Now().Unix())

}

func testAppointmentMadeNotifyClient(t *testing.T) {
	openId := "oDvnPw4zKAmraE2eccSUHinSya5E"
	examTime := "2020年07月29日"
	address := "又称南路253好"
	examCenterName := "美年大健康茂名店"
	AppointmentMadeNotifyClient(openId, examTime, examCenterName, address, 25)

}

func TestRefundLaunchedNotifyStaff(t *testing.T) {

	openIds := conf.C.RecvOpenIds
	RefundLaunchedNotifyStaff(openIds, "132564654564")
}

func TestRefundAgreedNotifyClient(t *testing.T) {
	RefundAgreedNotifyClient("oDvnPw4zKAmraE2eccSUHinSya5E", "7978789978", 98.65)
}
