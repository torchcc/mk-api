package service

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/wechat/oauth"
	. "mk-api/server/dao"
	"mk-api/server/model"
	"mk-api/server/util"
)

type WechatService interface {
	// UserExists(openId string) (bool, error)
	CheckUserNSetToken(resToken *oauth.ResAccessToken, oau *oauth.Oauth) (string, error)
}

type wechatService struct {
	model model.UserModel
}

func NewWechatService(userModel model.UserModel) WechatService {
	return &wechatService{
		model: userModel,
	}
}

func (service *wechatService) CheckUserNSetToken(resToken *oauth.ResAccessToken, oau *oauth.Oauth) (token string, err error) {
	// open_id.x123xua:{user_id: usr, mobile}
	// user_id_token.1232: token
	// token.xxx: {user_id: id, mobile: mobile}
	cli := Rdb.TokenRdbP.Get()
	defer cli.Close()

	openIdKey := "open_id" + resToken.OpenID

	// redis has openId-userInfo
	if res, _ := cli.Do("EXISTS", openIdKey); res.(int64) > 0 {

		userId, _ := redis.Int(cli.Do("HGET", openIdKey, "user_id"))
		mobile, _ := redis.String(cli.Do("HGET", openIdKey, "mobile"))
		return service.handleUserExists(userId, mobile, resToken, cli)

	}

	// mysql has openId-userInfo
	userId, mobile, err := service.model.FindUserByOpenId(resToken.OpenID)
	if err != nil {
		return service.handleUserExists(userId, mobile, resToken, cli)
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
		return "", err
	}
	// 设置 open_id.x123xua:{user_id: usr, mobile}
	service.setOpenIdUserInfo(openIdKey, userId, "", cli)

	// 设置 user_id_token.1232: token
	// 设置 token.xxx: {user_id: id, mobile: mobile}
	service.setToken(util.OpenId2Token(resToken.OpenID), "", userId, cli)
	return

}

func (service *wechatService) setToken(token string, mobile string, userId int, cli redis.Conn) {
	userIdTokenKey := "user_id_token." + strconv.Itoa(userId)

	_ = cli.Send("SETEX", userIdTokenKey, token, time.Second*7200)

	tokenUserInfoKey := "token." + token
	_ = cli.Send("HSET", tokenUserInfoKey, "user_id", userId)
	_ = cli.Send("HSET", tokenUserInfoKey, "mobile", mobile)
	_ = cli.Send("EXPIRE", tokenUserInfoKey, time.Second*7200)
	_ = cli.Flush()
}

func (service *wechatService) setOpenIdUserInfo(openIdKey string, userId int, mobile string, cli redis.Conn) {
	_ = cli.Send("HSET", openIdKey, "user_id", userId)
	_ = cli.Send("HSET", openIdKey, "user_id", mobile)
	_ = cli.Flush()
}

func (service *wechatService) handleUserExists(userId int, mobile string, resToken *oauth.ResAccessToken, cli redis.Conn) (token string, err error) {
	userIdTokenKey := "user_id_token." + strconv.Itoa(userId)
	if token, err = redis.String(cli.Do("GET", userIdTokenKey)); err != nil {
		// token过期了
		token = util.OpenId2Token(resToken.OpenID)
		service.setToken(token, mobile, userId, cli)
	}
	// token 没过期
	return token, nil

}
