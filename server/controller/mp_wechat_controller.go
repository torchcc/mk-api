package controller

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/menu"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"mk-api/library/ecode"
	"mk-api/server/conf"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/middleware"
	"mk-api/server/model"
	"mk-api/server/service"
	"mk-api/server/util"
	"mk-api/server/util/consts"
)

// wechat 路由注册
func WeChatRegister(router *gin.RouterGroup) {
	var (
		userModel        model.UserModel       = model.NewUserModel()
		wechatService    service.WechatService = service.NewWechatService(userModel)
		wechatController WeChatController      = NewWechatController(wechatService)
	)
	router.GET("/", wechatController.DockWithWeChatServer)
	router.POST("/", wechatController.WXMsgReceive)
	router.GET("/js_ticket", wechatController.JsApiTicket)
	router.GET("/enter_url", wechatController.GetEnterUrl)
	router.GET("/enter", wechatController.Enter)
	router.POST("/menu", wechatController.CreateMenu)
	router.GET("/menu", wechatController.ListMenu)
}

type WeChatController interface {
	DockWithWeChatServer(ctx *gin.Context)
	JsApiTicket(ctx *gin.Context)
	WXMsgReceive(ctx *gin.Context)
	GetEnterUrl(ctx *gin.Context)
	Enter(ctx *gin.Context)
	CreateMenu(ctx *gin.Context)
	ListMenu(ctx *gin.Context)
}

type wechatController struct {
	affAcc  *officialaccount.OfficialAccount
	service service.WechatService
}

// 接收微信消息
func (c *wechatController) WXMsgReceive(ctx *gin.Context) {
	var msg message.MixMessage
	err := ctx.ShouldBindXML(&msg)
	if err != nil {
		util.Log.Errorf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}
	util.Log.Debugf("[消息接收] - 收到消息, 消息类型为: %v, 消息内容为: %v\n", msg.MsgType, msg.Content)
	c.WXMsgReply(ctx, &msg)
}

func (c *wechatController) WXMsgReply(ctx *gin.Context, mixMessage *message.MixMessage) {
	// if mixMessage.MsgType == message.MsgTypeText {
	tcMsg := message.NewTransferCustomer("")
	tcMsg.MsgType = message.MsgTypeTransfer
	tcMsg.CreateTime = time.Now().Unix()
	tcMsg.FromUserName = mixMessage.ToUserName
	tcMsg.ToUserName = mixMessage.FromUserName
	msg, err := xml.Marshal(tcMsg)
	if err != nil {
		util.Log.Errorf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = ctx.Writer.Write(msg)
	// }

}

func (c *wechatController) ListMenu(ctx *gin.Context) {
	m := dao.AffAcc.GetMenu()
	menus, err := m.GetMenu()
	if err != nil {
		util.Log.Errorf("failed to get menu, err: [%s]", err.Error())
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}
	middleware.ResponseSuccess(ctx, menus)
}

func (c *wechatController) CreateMenu(ctx *gin.Context) {
	m := dao.AffAcc.GetMenu()
	buttons := []*menu.Button{&menu.Button{
		Type:       "view",
		Name:       "预约体检",
		Key:        "signUp4Exam",
		URL:        "https://www.mkhealth.club",
		MediaID:    "",
		AppID:      conf.C.WeChat.AppID,
		PagePath:   "",
		SubButtons: nil,
	}}

	err := m.SetMenu(buttons)
	if err != nil {
		util.Log.Errorf("failed to create menu, err: [%s]", err)
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}
	middleware.ResponseSuccess(ctx, "ok")
}

// DockWechat godoc
// @Summary 对接微信
// @Description 与微信服务器对接，此接口请忽略
// @Tags WechatTag
// @Accept json
// @Produce json
// @Param  timestamp query string true "时间戳"
// @Param  signature query string true "签名"
// @Param  nonce query string true "NONCE"
// @Param  echostr query string true "回响字符串"
// @Success 200 {object} string
// @Router / [get]

func (c *wechatController) DockWithWeChatServer(ctx *gin.Context) {
	timestamp, nonce, signatureIn := ctx.Query("timestamp"), ctx.Query("nonce"), ctx.Query("signature")
	signatureGen := makeSignature(timestamp, nonce)

	if signatureGen != signatureIn {
		util.Log.Infof("signatureGen != signatureIn signatureGen=%s,signatureIn=%s\n", signatureGen, signatureIn)
		ctx.String(http.StatusOK, "%s", "wrong")

	} else {
		// 如果请求来自于微信，则原样返回echostr参数内容 以上完成后，接入验证就会生效，开发者配置提交就会成功。
		echostr := ctx.Query("echostr")
		ctx.String(http.StatusOK, "%s", echostr)
	}
}

func makeSignature(timestamp string, nonce string) string {
	// 1. 将 plat_token、timestamp、nonce三个参数进行字典序排序
	sl := []string{conf.C.WeChat.Token, timestamp, nonce}
	sort.Strings(sl)
	// 2. 将三个参数字符串拼接成一个字符串进行sha1加密
	s := sha1.New()
	_, _ = io.WriteString(s, strings.Join(sl, ""))

	return fmt.Sprintf("%x", s.Sum(nil))
}

// JsApiTicket godoc
// @Summary 获取签名
// @Description 获取jsApiTicket 签名
// @Tags WechatTag
// @Accept json
// @Produce json
// @Param  uri query string true "传入需要的调用js-sdk的uri地址"
// @Success 200 {object} dto.JsApiTicketOutPut
// @Router /wx/js_ticket [get]
func (c *wechatController) JsApiTicket(ctx *gin.Context) {
	uri := ctx.Query("uri")
	if uri == "" {
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("缺少请求参数uri"))
		return
	}

	js := c.affAcc.GetJs()
	cfg, err := js.GetConfig(consts.UrlPrefix + uri)
	if err != nil {
		util.Log.Errorf("failed to get JsApiTicket, Param: %s, err: %v", uri, err)
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}
	middleware.ResponseSuccess(ctx, dto.JsApiTicketOutPut{Signature: cfg.Signature})
}

// Launch Oauth godoc
// @Summary 获取微信入口url
// @Description 获取微信入口url
// @Tags WechatTag
// @Param uri query string false "需要设置的button 入口，eg: '/index.html'"
// @Success 200 {object} middleware.Response{data=dto.GetEnterUrlOutput} "进入微信的url"
// @Router /wx/enter_url [get]
func (c *wechatController) GetEnterUrl(ctx *gin.Context) {
	uri := ctx.Query("uri")
	oau := c.affAcc.GetOauth()
	url, err := oau.GetRedirectURL(consts.UrlPrefix+"/wx/enter?uri="+uri, "snsapi_userinfo", "")
	if err != nil {
		util.Log.Errorf("fail to launch a oauth2 to wechat server: %v", err)
		return
	}
	middleware.ResponseSuccess(ctx, dto.GetEnterUrlOutput{Url: url})
}

// Enter  godoc
// @Summary 拿到code后的回调地址
// @Description 拿到code后的回调地址
// @Tags WechatTag
// @Param code query string true "微信服务器的code"
// @Success 200 {object} middleware.Response{data=dto.TokenOutput} "token"
// @Router /wx/enter [get]
func (c *wechatController) Enter(ctx *gin.Context) {
	oau := c.affAcc.GetOauth()
	code := ctx.Query("code")
	if code == "" {
		util.Log.Errorf("failed to get code from wechat server !")
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("failed to get code from wechat server"))
		return
	}
	util.Log.Debugf("微信的code是： %s", code)

	resToken, err := oau.GetUserAccessToken(code)
	if err != nil || resToken.OpenID == "" {
		errStr := fmt.Sprintf("failed to get access_token/open_id from wechat server, err: [%#v]", err)
		util.Log.Error(errStr)
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New(errStr))
		return
	}

	token, mobile, err := c.service.CheckUserNSetToken(&resToken, oau)

	if err != nil {
		util.Log.Errorf("查询用户设置token失败， open_id: %s, err: %v",
			resToken.OpenID, err)
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}
	util.Log.Debugf("handle enter logic done, token is [%s], mobile is [%s]", token, mobile)
	var mobileVerified int8 = 1
	if mobile == "" {
		mobileVerified = 0
	}
	middleware.ResponseSuccess(ctx, dto.TokenOutput{Token: token, MobileVerified: mobileVerified})
}

func NewWechatController(service service.WechatService) WeChatController {
	return &wechatController{
		affAcc:  dao.AffAcc,
		service: service,
	}
}
