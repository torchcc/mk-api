package service

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/wechat/oauth"
	"mk-api/library/ecode"
	. "mk-api/server/dao"
	"mk-api/server/model"
	"mk-api/server/util"
	tokenUtil "mk-api/server/util/token"
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

	openIdKey := "hash.open_id." + resToken.OpenID

	// redis has openId-userInfo
	if res, _ := cli.Do("EXISTS", openIdKey); res.(int64) > 0 {

		userId, _ := redis.Int64(cli.Do("HGET", openIdKey, "user_id"))
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
		return "", ecode.ServerErr
	}
	// 设置 open_id.x123xua:{user_id: usr, mobile}
	tokenUtil.SetOpenIdUserInfo(openIdKey, userId, "", cli)

	// 设置 user_id_token.1232: token
	// 设置 token.xxx: {user_id: id, mobile: mobile}
	tokenUtil.SetToken(tokenUtil.GenerateUuid(), "", userId, resToken.OpenID, cli)
	return
}

// 设置 user_id_token.1232: token
// 设置 token.xxx: {user_id: id, mobile: mobile}
func (service *wechatService) handleUserExists(userId int64, mobile string, resToken *oauth.ResAccessToken, cli redis.Conn) (token string, err error) {
	userIdTokenKey := "string.user_id_token." + strconv.FormatInt(userId, 10)
	token, err = redis.String(cli.Do("GET", userIdTokenKey))
	if err != nil {
		// token过期了
		token = tokenUtil.GenerateUuid()
		tokenUtil.SetToken(token, mobile, userId, resToken.OpenID, cli)
	}
	return
}
