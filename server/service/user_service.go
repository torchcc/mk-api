package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	"mk-api/server/dto"
	"mk-api/server/model"
	"mk-api/server/util"
	"mk-api/server/validator/id_card"
)

type UserService interface {
	Retrieve(id int64) (*dto.UserDetailOutput, error)
	FindAllAddrs(userId int64) (addrs []model.UserAddr, err error)
	SaveAddr(addr *model.UserAddr) (id int64, err error)
	RetrieveAddr(id int64) (addr *dto.GetUserAddrOutput, err error)
	DeleteAddr(id int64) (err error)
	UpdateUserAddr(ctx *gin.Context, id int64, addr *dto.UpdateUserAddrInput) (err error)

	FindAllExaminees(userId int64) ([]*dto.ListExamineeOutputEle, error)
	SaveExaminee(userId int64, input *dto.PostExamineeInput) (id int64, err error)
}

type userService struct {
	model         model.UserModel
	addrModel     model.UserAddrModel
	regionModel   model.RegionModel
	examineeModel model.ExamineeModel
}

func (service *userService) SaveExaminee(userId int64, input *dto.PostExamineeInput) (id int64, err error) {
	var bean = &dto.ExamineeBean{
		UserId:            userId,
		Gender:            0,
		CreateTime:        time.Now().Unix(),
		UpdateTime:        time.Now().Unix(),
		PostExamineeInput: dto.PostExamineeInput{},
	}
	tmp, _ := json.Marshal(input) // Marshal 可以传指针吗 test
	_ = json.Unmarshal(tmp, bean)
	if _, isMale, _, _ := id_card.GetCitizenNoInfo([]byte(input.IdCardNo)); isMale {
		bean.Gender = MALE
	} else {
		bean.Gender = Female
	}

	id, err = service.examineeModel.SaveExaminee(bean)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"user_id": userId}).Errorf("创建常用体检人失败, err: [%s]", err.Error())
	}
	return
}

func (service *userService) FindAllExaminees(userId int64) ([]*dto.ListExamineeOutputEle, error) {
	output, err := service.examineeModel.FindExamineesByUserId(userId)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"user_id": userId}).Errorf("获取常用体检人出错: [%s]", err.Error())
		return output, err
	}
	for _, examinee := range output {
		birthday, _, _, _ := id_card.GetCitizenNoInfo([]byte(examinee.IdCardNo))
		examinee.Age = (time.Now().Unix()-birthday)/(365*24*3600) + 1
	}
	return output, err
}

func (service *userService) SaveAddr(addr *model.UserAddr) (id int64, err error) {
	if addr.IsDefault == 1 {
		_ = service.addrModel.CancelOriginDefaultAddr(addr.UserId)
	}
	return service.addrModel.Save(addr)
}

func (service *userService) Retrieve(id int64) (*dto.UserDetailOutput, error) {
	u, err := service.model.FindUserByID(id)
	if err != nil {
		err = errors.Wrap(ecode.ServerErr,
			fmt.Sprintf("[FindUserByID] Params: [%v] failed with error: %s", id, err.Error()))
	}
	return u, err
}

func (service *userService) FindAllAddrs(userId int64) (addrs []model.UserAddr, err error) {
	addrs, err = service.addrModel.FindUserAddrByUserId(userId)
	if err != nil {
		util.Log.Errorf("查询用户收件地址列表出错, err: %s", err.Error())
		return
	}

	regionId2NameMap, err := service.regionModel.GetRegionIdNameMap()
	if err != nil {
		util.Log.Errorf("获取RegionId2NameMap单例出错, err: %s", err.Error())
		return

	}

	// 把id映射成省市区镇名称
	for i, _ := range addrs {
		addrs[i].ProvinceName = regionId2NameMap[addrs[i].ProvinceId]
		addrs[i].CityName = regionId2NameMap[addrs[i].CityId]
		addrs[i].CountyName = regionId2NameMap[addrs[i].CountyId]
		addrs[i].TownName = regionId2NameMap[addrs[i].TownId]
	}

	return
}

func (service *userService) RetrieveAddr(id int64) (addr *dto.GetUserAddrOutput, err error) {
	addr, err = service.addrModel.FindUserAddrByAddrId(id)
	if err != nil {
		util.Log.Errorf("查询用户收件地址出错, err: [%s]", err.Error())
	}
	return
}

func (service *userService) DeleteAddr(id int64) (err error) {
	err = service.addrModel.DeleteUserAddrByAddrId(id)
	if err != nil {
		util.Log.Errorf("删除用户收件地址出错, err: [%s]", err.Error())
	}
	return
}

func (service *userService) UpdateUserAddr(ctx *gin.Context, id int64, addr *dto.UpdateUserAddrInput) (err error) {
	if addr.IsDefault == 1 {
		_ = service.addrModel.CancelOriginDefaultAddr(ctx.GetInt64("userId"))
	}
	err = service.addrModel.UpdateUserAddr(id, addr)
	if err != nil {
		util.Log.Errorf("修改用户收件地址出错, err: [%s]", err.Error())
	}
	return
}

func NewUserService(userModel model.UserModel, addrModel model.UserAddrModel, regionModel model.RegionModel, examineeModel model.ExamineeModel) UserService {
	return &userService{
		model:         userModel,
		addrModel:     addrModel,
		regionModel:   regionModel,
		examineeModel: examineeModel,
	}
}
