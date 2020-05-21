package service

import (
	"github.com/sirupsen/logrus"
	"mk-api/server/model"
	"mk-api/server/util"
)

type UserService interface {
	Save(model.User) (uint32, error)
	Update(model.User) error
	Delete(model.User) error
	FindAll() ([]model.User, error)
	Retrieve(uint32) (*model.User, error)
}

type userService struct {
	model model.UserModel
}

func NewUserService(userModel model.UserModel) UserService {
	return &userService{
		model: userModel,
	}
}

func (service *userService) Save(user model.User) (uint32, error) { // user以后替换成dto 的对象， 不是do
	id, err := service.model.Save(user)
	return id, err
}

func (service *userService) Update(user model.User) error {
	return service.model.Update(user)
}

func (service *userService) Delete(user model.User) error {
	return service.model.Delete(user)
}

func (service *userService) FindAll() ([]model.User, error) {
	return service.model.FindAll()
}

func (service *userService) Retrieve(ID uint32) (*model.User, error) {
	u, err := service.model.FindUserByID(ID)
	if err != nil {
		util.Log.WithFields(logrus.Fields{"user_id": ID}).Errorf(
			"根据id查询用户信息失败: %s\n", err.Error())
	}
	return u, err
}
