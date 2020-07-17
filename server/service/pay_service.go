package service

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	wo "github.com/silenceper/wechat/v2/pay/order"
	"mk-api/library/ecode"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/pay/notify"
	"github.com/sirupsen/logrus"
	"mk-api/server/model"
	"mk-api/server/util"
	"mk-api/server/util/consts"
	wcUtil "mk-api/server/util/wechat"
)

type PayService interface {
	WechatPayCallBack(ctx *gin.Context) bool
	CheckPayStatus(ctx *gin.Context, prepayId string) (status int8, err error)
	Launch2ndPay(ctx *gin.Context, orderId int64) (*wo.Config, error)
}

type payService struct {
	payModel   model.PayModel
	orderModel model.OrderModel
	notify     *notify.Notify
}

func (service *payService) Launch2ndPay(ctx *gin.Context, orderId int64) (cfg *wo.Config, err error) {
	payStatus, err := service.orderModel.FindOrderPayStatusById(orderId)
	if err != nil {
		_ = ctx.Error(err)
		util.Log.WithFields(logrus.Fields{"order_id": orderId}).
			Errorf("failed to get order pay status, err: [%s]", err)
		return nil, ecode.NothingFound
	}
	if payStatus.TimeExpire < time.Now().Unix() || payStatus.Status != 0 {

		err = errors.New("该订单已经过期， 请重新下单")
		util.Log.WithFields(logrus.Fields{"order_id": orderId}).
			Warning(err.Error())
		_ = ctx.Error(err)
		return nil, ecode.RequestErr
	}
	cfg, err = wcUtil.Launch2ndPay(payStatus.NonceStr, payStatus.PrepayId)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"order_id": orderId}).
			Errorf("failed to calc paySign, err: [%s]", err.Error())
		_ = ctx.Error(err)
		return nil, ecode.ServerErr
	}
	return cfg, nil

}

func (service *payService) CheckPayStatus(ctx *gin.Context, prepayId string) (status int8, err error) {
	status, err = service.payModel.CheckPayStatusByPrepayId(prepayId)
	if err != nil {
		util.Log.Errorf("查询订单付款状态出错, err: [%s]", err)
	}
	return
}

func (service *payService) WechatPayCallBack(ctx *gin.Context) bool {
	var err error

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		util.Log.Errorf("read http body failed！err: [%s]", err.Error())
		return false
	}

	util.Log.Infof("wechat pay notify body: [%s]", string(body))

	var result notify.PaidResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		util.Log.Errorf("read http body xml failed! err: [%s]", err.Error())
		return false
	}

	if *result.ReturnCode == "FAIL" {
		util.Log.Errorf("notify result's return_code is FAIL, err: [%s]", *result.ReturnMsg)
		return false
	}

	// 这里加锁串行处理
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Lock()

	bill, err := service.payModel.FindBillByOutTradeNo(&result)
	if err != nil {
		util.Log.WithFields(
			logrus.Fields{"out_trade_no": result.OutTradeNo}).
			Errorf("微信notify查询原订单出错, err: [%s]", err.Error())
		return false
	}
	// 已经处理过， 直接返回SUCCESS
	if bill.Status == consts.Success && bill.TimeEnd != 0 && bill.TransactionId != "" {
		util.Log.WithFields(
			logrus.Fields{"out_trade_no": result.OutTradeNo}).
			Debug("微信notify, 已处理过该notify")
		return true
	}

	// 回调的订单总价与数据库价格不符
	if int(bill.TotalFee) != *result.TotalFee {
		util.Log.Warning(" total fee of notify result is not equal to the one in db")
		return false
	}

	// 进行签名校验
	if !service.notify.PaidVerifySign(result) {
		util.Log.Warning("notify result failed payVerifySign")
		return false
	}

	// 更新数据库
	if err = service.payModel.SuccessPaidResult2Bill(&result); err != nil {
		util.Log.Errorf("SuccessPaidResult2Bill failed, err: [%s]", err.Error())
		return false
	}
	if err = service.orderModel.UpdateOrderStatus(*result.OutTradeNo, consts.Success); err != nil {
		util.Log.Errorf("UpdateOrderStatus failed, err: [%s]", err.Error())
		return false

	}
	return true
}

func NewPayService(notify *notify.Notify, payModel model.PayModel, orderModel model.OrderModel) PayService {
	return &payService{
		payModel:   payModel,
		orderModel: orderModel,
		notify:     notify,
	}
}
