package service

import (
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
)

type RegionService interface {
	RetrieveRegionsByParentId(parentId int64) (output []*dto.Region, err error)
}

type regionService struct {
	regionModel model.RegionModel
}

func (service *regionService) RetrieveRegionsByParentId(parentId int64) (output []*dto.Region, err error) {

	output, err = service.regionModel.FindRegionsByParentId(parentId)
	if err != nil {
		util.Log.Errorf("list region failed, err: [%s]", err.Error())
	}
	return
}

func NewRegionService(regionModel model.RegionModel) RegionService {
	return &regionService{
		regionModel: regionModel,
	}
}
