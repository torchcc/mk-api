package wechat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount/message"
	"mk-api/server/dao"
	"mk-api/server/util"
)

const (
	orange     = "#642100"
	blue       = "#000093"
	softYellow = "#FFE4CA"
)

// 退款通知运营人员, 接受一个receiver list, 分别推送。
func RefundLaunchedNotifyStaff(openIds []string, outTradeNo string) {
	tmpl := dao.AffAcc.GetTemplate()
	const tmplId = "G5Rz2Ess-YYF6oomsw7JQQ5Et2GVHz5eOOxVJoCMEnY"
	curTime := time.Now().Format("2006年01月02日 15:04:05")

	for i := 0; i < len(openIds); i++ {
		msg := &message.TemplateMessage{
			ToUser:     openIds[i],
			TemplateID: tmplId,
			URL:        "",
			Color:      "",
			Data: map[string]*message.TemplateDataItem{
				"first": {
					Value: "平台有退款申请，请你进入后台完成审核",
				},
				"keyword1": {
					Value: outTradeNo,
					Color: blue,
				},
				"remark": {
					Value: fmt.Sprintf("请进入后台处理，申请时间 %s", curTime),
					Color: orange,
				},
			},
		}

		if _, err := tmpl.Send(msg); err != nil {
			util.Log.Errorf("failed to send msg to %s, err: [%s]", openIds[i], err.Error())
		}
	}
}

// 下单成功后微信推送给员工 use the 2nd
func OrderPaidNotifyStaff(openIds []string, outTradeNo string, amount float64, paidTime int64) {
	tmpl := dao.AffAcc.GetTemplate()
	const tmplId = "102fXlDTbJTx_RqhdLNh7KVZNJJOfbWo2AiwwtuA9A4"

	orderTimeStr := time.Unix(paidTime, 0).Format("2006年01月02日 15:04:05")
	for i := 0; i < len(openIds); i++ {
		msg := &message.TemplateMessage{
			ToUser:     openIds[i],
			TemplateID: tmplId,
			URL:        "",
			Color:      "",
			Data: map[string]*message.TemplateDataItem{
				"first": {
					Value: "您有新的订单，请在后台管理中心处理该订单",
					Color: "",
				},
				"keyword1": { // 订单号
					Value: outTradeNo,
					Color: blue,
				},
				"keyword2": { // 实付金额
					Value: strconv.FormatFloat(amount, 'f', -1, 64) + " 元",
					Color: orange,
				},
				"keyword3": { // 下单时间
					Value: orderTimeStr,
					Color: "",
				},
				"remark": { // 下单时间
					Value: "请尽快处理！",
					Color: "",
				},
			},
		}
		if _, err := tmpl.Send(msg); err != nil {
			util.Log.Errorf("failed to send msg to %s, err: [%s]", openIds[i], err.Error())
		}
	}
}

// 付款成功后推送给客户 give up the 3rd , use the 6th
func OrderPaidNotifyClient(openId, outTradeNo string, amount float64, orderId, paidTime int64) {
	tmpl := dao.AffAcc.GetTemplate()
	const tmplId = "jIWVI8mZj7C_v_PscxgHB1MslRApfe_yE0q1ScXQgZ0"
	orderTimeStr := time.Unix(paidTime, 0).Format("2006年01月02日 15:04:05")
	url := fmt.Sprintf("https://www.mkhealth.club/#/orderDetail?orderNum=%d&state=2", orderId)
	msg := &message.TemplateMessage{
		ToUser:     openId,
		TemplateID: tmplId,
		URL:        url,
		Color:      "",
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: "您好，您的体检套餐已成功下单!",
				Color: "",
			},
			"keyword1": { // 订单编号
				Value: outTradeNo,
				Color: blue,
			},
			"keyword2": { // 下单时间
				Value: orderTimeStr,
				Color: "",
			},
			"keyword3": { // 订单类型
				Value: "体检",
				Color: "",
			},
			"keyword4": { // 订单金额
				Value: strconv.FormatFloat(amount, 'f', 2, 64) + " 元",
				Color: "",
			},
			"remark": {
				Value: "温馨提示：您的订单信息客服会在一个工作日内与您电话核实，请您保持电话通畅。如有疑问请拨打客服热线0668-2853837。祝您身体健康，幸福美满！",
				Color: orange,
			},
		},
	}
	if _, err := tmpl.Send(msg); err != nil {
		util.Log.Errorf("failed to send msg to %s, err: [%s]", openId, err.Error())
	}

}

// 人工客服预约成功后推送給客户 admin 用的
func AppointmentMadeNotifyClient(openId string, examTime string, examCenterName string, address string, orderId int64) {
	tmpl := dao.AffAcc.GetTemplate()
	const tmplId = "-XUWF_622novQ6keJug8MEWjpuqIrBcw6H7Yvc21CPs"
	url := fmt.Sprintf("https://www.mkhealth.club/#/orderDetail?orderNum=%d&state=2", orderId)
	msg := &message.TemplateMessage{
		ToUser:     openId,
		TemplateID: tmplId,
		URL:        url,
		Color:      "",
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: "您好！您的体检预约成功。",
				Color: "",
			},
			"keyword1": { // 体检时间
				Value: examTime,
				Color: blue,
			},
			"keyword2": { // 医院名称， 门店
				Value: examCenterName,
				Color: "",
			},
			"keyword3": { // 地址
				Value: address,
				Color: "",
			},
			"keyword4": { // 抽血时间
				Value: "7:30-9:30",
				Color: blue,
			},
			"remark": {
				Value: "如有疑问请拨打客服热线0668-2853837。祝您身体健康，幸福美满！",
				Color: orange,
			},
		},
	}

	if _, err := tmpl.Send(msg); err != nil {
		util.Log.Errorf("failed to send msg to %s, err: [%s]", openId, err.Error())
	}
}

// 人工审核退款通过后， 推送给客户， admin用
func RefundAgreedNotifyClient(openId, outTradeNo string, amount float64) {
	tmpl := dao.AffAcc.GetTemplate()
	const tmplId = "De7WxIRy_ke0PiadqQjcUpIpHo1GQCa9gNVyr7zCp9A"
	msg := &message.TemplateMessage{
		ToUser:     openId,
		TemplateID: tmplId,
		URL:        "",
		Color:      "",
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: "商家已同意您的退款申请",
				Color: "",
			},
			"keyword1": { // 订单编号
				Value: outTradeNo,
				Color: blue,
			},
			"keyword2": { // 退款状态
				Value: "退款审核通过",
				Color: "",
			},
			"keyword3": { // 退款金额
				Value: "现金" + strconv.FormatFloat(amount, 'f', -1, 64) + "元",
				Color: blue,
			},
			"remark": { // 备注
				Value: "商家已同意您的退款申请，系统会在2-3天内提交微信处理，微信审核后在2-5个工作日自动原路退回至您的支付账户。若超时未收到退款，请联系官方客服核实",
				Color: "",
			},
		},
	}
	if _, err := tmpl.Send(msg); err != nil {
		util.Log.Errorf("failed to send msg to %s, err: [%s]", openId, err.Error())
	}
}
