package service

import (
	"sort"

	"github.com/gin-gonic/gin"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
)

type PkgAttr = int8

const (
	ITEM      PkgAttr = 1
	NOTICE            = 2
	PROCEDURE         = 3
)

type PackageService interface {
	ListPackage(ctx *gin.Context, input *dto.ListPackageInput) (data *dto.PaginateListOutput, err error)
	RetrievePackage(id int64) (data *dto.GetPackageOutPut, err error)
	ListDisease() ([]*dto.Disease, error)
	ListCategory() ([]*dto.Category, error)
}

type packageService struct {
	packageModel model.PackageModel
}

func (service *packageService) ListCategory() ([]*dto.Category, error) {
	return service.packageModel.ListCategory()
}

func (service *packageService) ListDisease() ([]*dto.Disease, error) {
	return service.packageModel.ListDisease()
}

func (service *packageService) RetrievePackage(id int64) (*dto.GetPackageOutPut, error) {
	var data dto.GetPackageOutPut

	basicInfo, err := service.packageModel.FindPackageBasicInfo(id)
	if err != nil {
		util.Log.Errorf("获取套餐基本信息出错, id: [%d], err: [%s]", id, err.Error())
		return &data, err
	}
	data.BasicInfo = basicInfo

	attrs, err := service.packageModel.FindPackageAttr(id)
	if err != nil {
		util.Log.Errorf("获取套餐属性出错, id: [%d], err: [%s]", id, err.Error())
		return &data, err
	}

	var (
		items     []dto.PackageItem
		notices   []dto.PackageNotice
		procedure []dto.PackageProcedure
	)
	for _, attr := range attrs {
		switch attr.AttrType {
		case ITEM:
			items = append(items, attr)
		case NOTICE:
			notices = append(notices, attr)
		case PROCEDURE:
			procedure = append(procedure, attr)
		}
	}

	sort.Sort(dto.PackageAttributes(items))
	sort.Sort(dto.PackageAttributes(notices))
	sort.Sort(dto.PackageAttributes(procedure))
	data.Items = items
	data.Notices = notices
	data.Procedure = procedure
	return &data, err
}

func (service *packageService) ListPackage(ctx *gin.Context, input *dto.ListPackageInput) (*dto.PaginateListOutput, error) {
	var data dto.PaginateListOutput
	list, err := service.packageModel.ListPackage(input)
	if err != nil {
		util.Log.Errorf("查询套餐列表出错, err: [%s]", err.Error())
		return &data, err
	}
	var length int
	if len(list) == int(input.PageSize)+1 {
		data.HasNext = 1
		length = len(list) - 1
	} else {
		length = len(list)
	}
	list = list[:length]
	data.PageSize = int64(len(list))
	data.PageNo = input.PageNo
	data.List = list
	return &data, err
}

func NewPackageService(packageModel model.PackageModel) PackageService {
	return &packageService{
		packageModel: packageModel,
	}
}
