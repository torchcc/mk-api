package service

import (
	"errors"
	"sort"

	"github.com/gin-gonic/gin"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
)

type CargoService interface {
	RetrieveCargo(ctx *gin.Context) ([]dto.GetCargoOutputElem, error)
	CreateCargo(ctx *gin.Context, pkgId int64) (err error)
	RemoveCargoEntries(input *dto.DeleteCargoEntriesInput) (err error)
}

type cargoService struct {
	cargoModel   model.CargoModel
	packageModel model.PackageModel
}

func (service *cargoService) RemoveCargoEntries(input *dto.DeleteCargoEntriesInput) (err error) {
	return service.cargoModel.RemoveCargoEntries(input)
}

func (service *cargoService) CreateCargo(ctx *gin.Context, pkgId int64) (err error) {
	userId := ctx.GetInt64("userId")
	if id := service.cargoModel.FindCargoItemId(userId, pkgId); id != 0 {
		err = service.cargoModel.IncrementPkgCount(id)
		return
	}
	price, err := service.packageModel.FindPackagePriceById(pkgId)
	if err != nil || price == 0.0 {
		util.Log.Errorf("套餐不存在，pkg_id: [%v], err: [%v]", pkgId, err)
		err = errors.New("套餐不存在")
	} else {
		err = service.cargoModel.CreateCargo(userId, pkgId)
	}
	return
}

func (service *cargoService) RetrieveCargo(ctx *gin.Context) (pkgs []dto.GetCargoOutputElem, err error) {
	userId := ctx.GetInt64("userId")
	pkgs, err = service.cargoModel.FindCargoByUserId(userId)
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

func NewCargoService(cargoModel model.CargoModel, packageModel model.PackageModel) CargoService {
	return &cargoService{
		cargoModel:   cargoModel,
		packageModel: packageModel,
	}
}
