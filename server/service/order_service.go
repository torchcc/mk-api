package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/pay"
	wo "github.com/silenceper/wechat/v2/pay/order"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/server/conf"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
	"mk-api/server/util/consts"
	"mk-api/server/util/token"
	"mk-api/server/util/xtime"
	"mk-api/server/validator/id_card"
)

type OrderService interface {
	CreateOrder(ctx *gin.Context, input *dto.PostOrderInput) (*wo.Config, error)
	ListOrder(ctx *gin.Context, input *dto.ListOrderInput) (*dto.PaginateListOutput, error)
	RetrieveOrder(ctx *gin.Context, id int64) (*dto.RetrieveOrderOutput, error)
	RemoveOrder(ctx *gin.Context, id int64) error
}

type orderService struct {
	orderModel   model.OrderModel
	cartModel    model.CartModel
	packageModel model.PackageModel
	payModel     model.PayModel
	wechatPay    *pay.Pay
}

func (service *orderService) RemoveOrder(ctx *gin.Context, id int64) error {
	userId := ctx.GetInt64("userId")
	err := service.orderModel.DeleteOrderByIdNUserId(userId, id)
	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"user_id":  userId,
			"order_id": id,
		}).Errorf("删除订单失败， err: [%s]", err)
	}
	return err
}

func (service *orderService) RetrieveOrder(ctx *gin.Context, id int64) (*dto.RetrieveOrderOutput, error) {
	output, err := service.orderModel.FindOrderDetailById(id, service.packageModel)
	if err != nil {
		util.Log.Errorf("获取订单详情出错, err: [%]", err)
	}
	return output, err
}

func (service *orderService) ListOrder(ctx *gin.Context, input *dto.ListOrderInput) (*dto.PaginateListOutput, error) {
	var data dto.PaginateListOutput
	list, err := service.orderModel.ListOrder(input, ctx.GetInt64("userId"))
	if err != nil {
		util.Log.Errorf("查询订单列表出错, err: [%s]", err)
		return &data, err
	}
	data.PageSize = int64(len(list))
	if data.PageSize == input.PageSize+1 {
		data.HasNext = 1
	}
	data.PageNo = input.PageNo
	data.List = list
	return &data, err
}

func (service *orderService) CreateOrder(ctx *gin.Context, input *dto.PostOrderInput) (*wo.Config, error) {
	var err error

	userId := ctx.GetInt64("userId")
	orderItems := make([]*dto.OrderItem, 0, 4)
	cartIds := make([]int64, 0, 8)

	var amount float64
	for _, cItem := range input.CartItems {
		cartIds = append(cartIds, cItem.CartId)

		// 检查套餐是否存在
		priceNTargetInfo, err := service.packageModel.FindPackagePriceNTargetById(cItem.PackageId)
		if err != nil {
			_ = ctx.Error(err)
			util.Log.WithFields(logrus.Fields{"userId": userId}).Errorf("套餐不存在: [%v]", input)
			return nil, ecode.ServerErr
		}
		// 检查套餐数量和体检人数量
		diff := cItem.PackageCount - len(cItem.Examinees)
		if diff < 0 {
			cItem.Examinees = cItem.Examinees[:len(cItem.Examinees)]
			diff = 0
		}
		// 创建orderItem 对象
		for i := 0; i < len(cItem.Examinees); i++ {
			// 鉴定target 和选择性别
			if priceNTargetInfo.Target != AnyGender {
				_, isMale, _, _ := id_card.GetCitizenNoInfo([]byte(cItem.Examinees[i].IdCardNo))
				if (priceNTargetInfo.Target == MALE) != isMale {
					var errStr string
					switch priceNTargetInfo.Target {
					case MALE:
						errStr = `此为'男性'套餐，女性人员是无法体检的，请悉知`
					case UnMarriedFemale, MarriedFemale:
						errStr = `此为'女性'套餐，男性人员是无法体检的，请悉知`
					}
					_ = ctx.Error(errors.New(errStr))
					return nil, ecode.RequestErr
				}
			}

			// 鉴定体检日期
			tomorrow := xtime.TomorrowStartAt()
			if cItem.Examinees[i].ExamineDate < tomorrow {
				_ = ctx.Error(errors.New("需至少提前一天预约体检"))
				return nil, ecode.RequestErr
			}

			orderItem := &dto.OrderItem{
				UserId:       userId,
				OrderId:      0,
				PackageId:    cItem.PackageId,
				PackagePrice: priceNTargetInfo.Price,
				CreateTime:   time.Now().Unix(),
				UpdateTime:   time.Now().Unix(),
				Examinee:     cItem.Examinees[i],
			}
			orderItems = append(orderItems, orderItem)
			amount += priceNTargetInfo.Price
		}

		for i := 0; i < diff; i++ {
			orderItem := &dto.OrderItem{
				UserId:     userId,
				OrderId:    0,
				PackageId:  cItem.PackageId,
				CreateTime: time.Now().Unix(),
				UpdateTime: time.Now().Unix(),
				Examinee:   &dto.Examinee{},
			}
			orderItems = append(orderItems, orderItem)
			amount += priceNTargetInfo.Price
		}

	}

	// 雪花算法产生 outTradeNo
	outTradeNo, err := token.GenerateSnowflake()
	if err != nil {
		util.Log.WithFields(logrus.Fields{"userId": userId}).Errorf("雪花算法出错， 请求体信息: [%v]", input)
		return nil, ecode.ServerErr
	}

	order := dto.Order{
		OutTradeNo: outTradeNo.Int64(),
		UserId:     userId,
		Mobile:     input.SubscriberMobile,
		OpenId:     ctx.GetString("openId"),
		Amount:     amount,
		Remark:     input.SubscriberComment,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	order.Id, err = service.orderModel.SaveOrder(&order, orderItems)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"userId": userId}).Errorf("创建订单出错, 请求参数: [%v]", input)
		return nil, ecode.ServerErr
	}

	cfg, err := service.makeWechatOrderNPrepay(ctx, &order)
	if err != nil {
		return nil, ecode.ServerErr
	}

	// 最后生成预付单后才删除购物车
	if err = service.cartModel.RemoveCartEntries(cartIds); err != nil {
		util.Log.Errorf("更新购物车条目出错, err: [%s]", err)
	}
	return cfg, nil
}

func (service *orderService) makeWechatOrderNPrepay(ctx *gin.Context, order *dto.Order) (*wo.Config, error) {
	var err error

	util.Log.WithFields(logrus.Fields{
		"user_id": ctx.GetInt64("userId"),
	}).Infof("用户的IP: [%s]", ctx.ClientIP())

	// 微信统一下单
	params := &wo.Params{
		TotalFee:   strconv.Itoa(int(order.Amount)),
		CreateIP:   ctx.ClientIP(),
		Body:       "迈康-体检套餐",
		OutTradeNo: strconv.FormatInt(order.OutTradeNo, 10),
		OpenID:     ctx.GetString("openId"),
		TradeType:  "JSAPI",
		SignType:   "MD5",
		Detail:     "预约体检套餐",
		Attach:     "迈康体检",
		GoodsTag:   "",
		NotifyURL:  conf.C.WeChat.PayNotifyURL,
	}
	wechatOrder := service.wechatPay.GetOrder() // 获取微信订单对象

	cfg, err := wechatOrder.BridgeConfig(params) // 下单+获取prepayId+获取返回给前端的cfg
	if err != nil {
		util.Log.Errorf("调用微信统一下单出错, err: [%s]", err)
		return nil, err
	}

	// 创建 mkb_trade_bill 条目
	now := time.Now().Unix()
	timeExpire := now + 3600*2

	bill := &dto.TradeBill{
		OrderId:    order.Id,
		OutTradeNo: order.OutTradeNo,
		PrepayId:   cfg.PrePayID,
		NonceStr:   cfg.NonceStr,
		TotalFee:   int64(order.Amount),
		FeeType:    Income,
		Status:     0,
		TransType:  Earned,
		TimeStart:  now,
		TimeExpire: timeExpire,
		CreateTime: now,
		UpdateTime: now,
	}

	billId, err := service.payModel.SaveTradeBill(bill)

	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"order_id": order.Id,
		}).Errorf("生成支付流水出错, err: [%s]", err.Error())
		return nil, err
	}

	util.Log.WithFields(logrus.Fields{
		"order_id": order.Id,
		"bill_id":  billId,
	}).Infof("生成支付流水成功!")

	time.AfterFunc(consts.OrderExpireIn, func() {
		_ = service.payModel.ExpireBill(billId)
		_ = service.orderModel.UpdateOrderStatus(params.OutTradeNo, consts.Closed)
	})

	return &cfg, nil
}

func NewOrderService(orderModel model.OrderModel, packageModel model.PackageModel,
	cartModel model.CartModel, payModel model.PayModel, wechatPay *pay.Pay) OrderService {
	return &orderService{
		orderModel:   orderModel,
		packageModel: packageModel,
		cartModel:    cartModel,
		wechatPay:    wechatPay,
		payModel:     payModel,
	}
}
