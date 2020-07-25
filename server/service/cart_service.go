package service

import (
	"errors"
	"sort"

	"github.com/gin-gonic/gin"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
)

type CartService interface {
	RetrieveCart(ctx *gin.Context) ([]dto.GetCartOutputElem, error)
	CreateCart(ctx *gin.Context, pkgId int64, pkgCount int64) (err error)
	RemoveCartEntries(input *dto.DeleteCartEntriesInput) (err error)
}

type cartService struct {
	cartModel    model.CartModel
	packageModel model.PackageModel
}

func (service *cartService) RemoveCartEntries(input *dto.DeleteCartEntriesInput) (err error) {
	return service.cartModel.RemoveCartEntries(input.CartIds)
}

func (service *cartService) CreateCart(ctx *gin.Context, pkgId int64, pkgCount int64) (err error) {
	userId := ctx.GetInt64("userId")
	if id := service.cartModel.FindCartItemId(userId, pkgId); id != 0 {
		err = service.cartModel.IncrementPkgCount(id, pkgCount)
		return
	}
	_, err = service.packageModel.FindPackagePriceNTargetById(pkgId)
	if err != nil {
		util.Log.Errorf("套餐不存在，pkg_id: [%v], err: [%v]", pkgId, err)
		err = errors.New("套餐不存在")
	} else {
		err = service.cartModel.CreateCart(userId, pkgId, pkgCount)
	}
	return
}

func (service *cartService) RetrieveCart(ctx *gin.Context) (pkgs []dto.GetCartOutputElem, err error) {
	userId := ctx.GetInt64("userId")
	pkgs, err = service.cartModel.FindCartByUserId(userId)
	if err != nil {
		util.Log.Errorf("查找用户购物车失败, userId: [%d], err: [%s]", userId, err.Error())
		return
	}
	// 按更新时间递减排序
	if len(pkgs) > 0 {
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].UpdateTime > pkgs[j].UpdateTime
		})
	}
	return
}

func NewCartService(cartModel model.CartModel, packageModel model.PackageModel) CartService {
	return &cartService{
		cartModel:    cartModel,
		packageModel: packageModel,
	}
}
