package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
	"mk-api/server/util/token"
	"mk-api/server/util/xtime"
	"mk-api/server/validator/id_card"
)

type PackageTarget = int8

const (
	AnyGender       PackageTarget = 0
	MALE                          = 1
	UnMarriedFemale               = 2
	MarriedFemale                 = 3
)

type OrderService interface {
	CreateOrder(ctx *gin.Context, input *dto.PostOrderInput) (err error)
}

type orderService struct {
	orderModel   model.OrderModel
	cartModel    model.CartModel
	packageModel model.PackageModel
}

func (service *orderService) CreateOrder(ctx *gin.Context, input *dto.PostOrderInput) (err error) {
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
			return ecode.ServerErr
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
					return ecode.RequestErr
				}
			}

			// 鉴定体检日期
			tomorrow := xtime.TomorrowStartAt()
			if cItem.Examinees[i].ExamineDate < tomorrow {
				_ = ctx.Error(errors.New("需至少提前一天预约体检"))
				return ecode.RequestErr
			}

			orderItem := &dto.OrderItem{
				UserId:     userId,
				OrderId:    0,
				PackageId:  cItem.PackageId,
				CreateTime: time.Now().Unix(),
				UpdateTime: time.Now().Unix(),
				Examinee:   cItem.Examinees[i],
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

	// 雪花算法产生 biz_no
	snowflake, err := token.GenerateSnowflake()
	if err != nil {
		util.Log.WithFields(logrus.Fields{"userId": userId}).Errorf("雪花算法出错， 请求体信息: [%v]", input)
		return ecode.ServerErr
	}

	order := dto.Order{
		BizNo:      snowflake,
		UserId:     userId,
		Mobile:     input.SubscriberMobile,
		OpenId:     ctx.GetString("openId"),
		Amount:     amount,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	orderId, err := service.orderModel.SaveOrder(&order, orderItems)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"userId": userId}).Errorf("创建订单出错, 请求参数: [%v]", input)
		return ecode.ServerErr
	}
	if err = service.cartModel.RemoveCartEntries(cartIds); err != nil {
		util.Log.Errorf("删除购物车条目出错, err: [%s]", err)
	}
	// TODO 到这人了 下一步发起微信支付。
	fmt.Println(orderId)
	return nil
}

func NewOrderService(orderModel model.OrderModel, packageModel model.PackageModel, cartModel model.CartModel) OrderService {
	return &orderService{
		orderModel:   orderModel,
		packageModel: packageModel,
		cartModel:    cartModel,
	}
}
