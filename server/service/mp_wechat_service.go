package service

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/sirupsen/logrus"
	"mk-api/library/ecode"
	. "mk-api/server/dao"
	"mk-api/server/model"
	"mk-api/server/util"
	tokenUtil "mk-api/server/util/token"
)

type WechatService interface {
	// UserExists(openId string) (bool, error)
	CheckUserNSetToken(resToken *oauth.ResAccessToken, oau *oauth.Oauth) (token string, retMobile string, err error)
}

type wechatService struct {
	model model.UserModel
}

func NewWechatService(userModel model.UserModel) WechatService {
	return &wechatService{
		model: userModel,
	}
}

func (service *wechatService) CheckUserNSetToken(resToken *oauth.ResAccessToken, oau *oauth.Oauth) (token string, retMobile string, err error) {
	// open_id.x123xua:{user_id: usr, mobile}
	// user_id_token.1232: token
	// token.xxx: {user_id: id, mobile: mobile}
	cli := Rdb.TokenRdbP.Get()
	defer cli.Close()

	openIdKey := "hash.open_id." + resToken.OpenID

	// redis has openId-userInfo
	if res, _ := cli.Do("EXISTS", openIdKey); res.(int64) > 0 {
		util.Log.Debugf("该用户存在， open_id_key 为 [%s]", openIdKey)
		userId, _ := redis.Int64(cli.Do("HGET", openIdKey, "user_id"))
		mobile, _ := redis.String(cli.Do("HGET", openIdKey, "mobile"))
		token, err = service.handleUserExists(userId, mobile, resToken, cli)
		return token, mobile, err
	}

	// mysql has openId-userInfo
	userId, mobile, err := service.model.FindUserByOpenId(resToken.OpenID)
	if err == nil {
		token, err = service.handleUserExists(userId, mobile, resToken, cli)
		return token, mobile, err
	}

	// user Does not exists
	wechatUserInfo, err := oau.GetUserInfo(resToken.AccessToken, resToken.OpenID)
	var u model.User
	u.OpenId = resToken.OpenID
	if err != nil {
		util.Log.Errorf("拉取微信用户信息错误: %v", err)
	} else {
		u.UserName = wechatUserInfo.Nickname
		u.AvatarUrl = wechatUserInfo.HeadImgURL
		u.Gender = wechatUserInfo.Sex
		u.Country = wechatUserInfo.Country
		u.Province = wechatUserInfo.Province
		u.City = wechatUserInfo.City
	}
	userId, err = service.model.Save(&u)
	if err != nil {
		util.Log.Errorf("创建用户失败: err: %v, openId: %s", err, resToken.OpenID)
		return "", "", ecode.ServerErr
	}
	// 设置 open_id.x123xua:{user_id: usr, mobile}
	tokenUtil.SetOpenIdUserInfo(openIdKey, userId, "", cli)

	// 设置 user_id_token.1232: token
	// 设置 token.xxx: {user_id: id, mobile: mobile}
	token = tokenUtil.GenerateUuid()
	tokenUtil.SetToken(token, "", userId, resToken.OpenID, cli)
	return token, "", nil
}

// 设置 user_id_token.1232: token
// 设置 token.xxx: {user_id: id, mobile: mobile}
func (service *wechatService) handleUserExists(userId int64, mobile string, resToken *oauth.ResAccessToken, cli redis.Conn) (token string, err error) {
	userIdTokenKey := "string.user_id_token." + strconv.FormatInt(userId, 10)
	token, err = redis.String(cli.Do("GET", userIdTokenKey))
	if err != nil {
		// token过期了
		util.Log.WithFields(logrus.Fields{
			"user_id": userId,
			"mobile":  mobile,
		}).Infof("token expired, user_id is: [%d], mobile is: [%s]", userId, mobile)
		token = tokenUtil.GenerateUuid()
		tokenUtil.SetToken(token, mobile, userId, resToken.OpenID, cli)
	}
	return token, nil
}
