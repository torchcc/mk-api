package service

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	. "mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
	"mk-api/server/util/consts"
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
	var (
		ctgs, cacheCtgs []*dto.Category
	)
	key := consts.CacheCategory
	if Rdb.ApiCache.Exists(key) {
		if data, err := Rdb.ApiCache.Get(key); err != nil {
			util.Log.Errorf("failed to get categories from redis, err: %s", err.Error())
		} else {
			_ = json.Unmarshal(data, &cacheCtgs)
			util.Log.Debugf("hit redis when getting category list!")
			return cacheCtgs, nil
		}
	}
	ctgs, err := service.packageModel.ListCategory()
	if err != nil {
		return nil, err
	}
	go Rdb.ApiCache.SetEx(key, ctgs, consts.CategoryListDuration)
	return ctgs, nil

}

func (service *packageService) ListDisease() ([]*dto.Disease, error) {
	var (
		diseases, cacheDiseases []*dto.Disease
	)
	key := consts.CacheDisease
	if Rdb.ApiCache.Exists(key) {
		if data, err := Rdb.ApiCache.Get(key); err != nil {
			util.Log.Warningf("failed to get disease from redis, err: %s", err.Error())
		} else {
			_ = json.Unmarshal(data, &cacheDiseases)
			util.Log.Debugf("hit redis when getting disease list!")
			return cacheDiseases, nil
		}
	}
	diseases, err := service.packageModel.ListDisease()
	if err != nil {
		return nil, err
	}
	go Rdb.ApiCache.SetEx(key, diseases, consts.DiseaseListDuration)
	return diseases, nil
}

func (service *packageService) RetrievePackage(id int64) (*dto.GetPackageOutPut, error) {
	var output dto.GetPackageOutPut

	// try to get result from redis
	key := consts.CachePackage + "." + strconv.FormatInt(id, 10)
	if Rdb.ApiCache.Exists(key) {
		data, err := Rdb.ApiCache.Get(key)
		if err != nil {
			util.Log.Warningf("failed to retrieve package from redis, err: %s", err.Error())
		} else {
			util.Log.Debugf("hit redis when retrieving package !")
			_ = json.Unmarshal(data, &output)
			return &output, nil
		}
	}

	basicInfo, err := service.packageModel.FindPackageBasicInfo(id)
	if err != nil {
		util.Log.Errorf("获取套餐基本信息出错, id: [%d], err: [%s]", id, err.Error())
		return &output, err
	}
	output.BasicInfo = basicInfo

	attrs, err := service.packageModel.FindPackageAttr(id)
	if err != nil {
		util.Log.Errorf("获取套餐属性出错, id: [%d], err: [%s]", id, err.Error())
		return &output, err
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
	output.Items = items
	output.Notices = notices
	output.Procedure = procedure

	go Rdb.ApiCache.SetEx(key, output, consts.PackageOneDuration)

	return &output, err
}

func (service *packageService) ListPackage(ctx *gin.Context, input *dto.ListPackageInput) (*dto.PaginateListOutput, error) {
	var (
		output, cacheOutput dto.PaginateListOutput
	)
	key := input.GetListKey()
	if Rdb.ApiCache.Exists(key) {
		data, err := Rdb.ApiCache.Get(key)
		if err != nil {
			util.Log.Warningf("failed to pkg list from redis, err: %s", err.Error())
		} else {
			util.Log.Debugf("hit redis when getting package list !")
			_ = json.Unmarshal(data, &cacheOutput)
			return &cacheOutput, nil
		}
	}

	list, err := service.packageModel.ListPackage(input)
	if err != nil {
		util.Log.Errorf("查询套餐列表出错, err: [%s]", err.Error())
		return &output, err
	}
	var length int
	if len(list) == int(input.PageSize)+1 {
		output.HasNext = 1
		length = len(list) - 1
	} else {
		length = len(list)
	}
	list = list[:length]
	output.PageSize = int64(length)
	output.PageNo = input.PageNo
	output.List = list

	go Rdb.ApiCache.SetEx(key, output, consts.PackageListDuration)

	return &output, err
}

func NewPackageService(packageModel model.PackageModel) PackageService {
	return &packageService{
		packageModel: packageModel,
	}
}
