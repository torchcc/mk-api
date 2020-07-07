package controller

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/message"
	"mk-api/library/ecode"
	"mk-api/server/conf"
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
	router.GET("/js_ticket", wechatController.JsApiTicket)
	router.GET("/launch_auth", wechatController.LaunchAuth)
	router.GET("/enter", wechatController.Enter)

}

type WeChatController interface {
	DockWithWeChatServer(ctx *gin.Context)
	JsApiTicket(ctx *gin.Context)
	Echo(ctx *gin.Context)
	LaunchAuth(ctx *gin.Context)
	Enter(ctx *gin.Context)
}

type wechatController struct {
	cfg     *wechat.Config
	wc      *wechat.Wechat
	service service.WechatService
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
// @Param  uri query string true "传入需要的调用js-sdk的uri地址" example "/static/index.html"
// @Success 200 {object} dto.JsApiTicketOutPut
// @Router /wx/js_ticket [get]
func (c *wechatController) JsApiTicket(ctx *gin.Context) {
	uri := ctx.Query("uri")
	if uri == "" {
		middleware.ResponseError(ctx, ecode.RequestErr, errors.New("缺少请求参数uri"))
		return
	}

	js := c.wc.GetJs()
	cfg, err := js.GetConfig(consts.UrlPrefix + uri)
	if err != nil {
		util.Log.Errorf("failed to get JsApiTicket, Param: %s, err: %v", uri, err)
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}
	middleware.ResponseSuccess(ctx, dto.JsApiTicketOutPut{Signature: cfg.Signature})
}

// Launch Oauth godoc
// @Summary 发起授权
// @Description 发起授权，直接在一个button中发一个get请求到这里即可
// @Tags WechatTag
// @Param uri query string false "需要设置的button 入口" example "/index.html"
// @Success 200 {object} string ""
// @Router /wx/launch_auth [get]
func (c *wechatController) LaunchAuth(ctx *gin.Context) {
	uri := ctx.Query("uri")
	oau := c.wc.GetOauth()
	url, err := oau.GetRedirectURL(consts.UrlPrefix+"/wx/enter?uri="+uri, "snsapi_userinfo", "")
	if err != nil {
		util.Log.Errorf("fail to launch a oauth2 to wechat server: %v", err)
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// TODO 待测试
// 发送请求到LaunchAuth 的时候， 重定向到Enter， Enter 重定向到index.html 或者其他页面
func (c *wechatController) Enter(ctx *gin.Context) {
	oau := c.wc.GetOauth()
	code := ctx.Query("code")
	if code == "" {
		util.Log.Errorf("failed to get code from wechat server !")
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("failed to get code from wechat server"))
		return
	}

	resToken, err := oau.GetUserAccessToken(code)
	if err != nil || resToken.OpenID == "" {
		util.Log.Errorf("failed to get access_token/open_id from wechat server !")
		middleware.ResponseError(ctx, ecode.ServerErr, errors.New("failed to get access_token/open_id from wechat server"))
		return
	}

	token, err := c.service.CheckUserNSetToken(&resToken, oau)

	if err != nil {
		util.Log.Errorf("查询用户设置token失败， open_id: %s, err: %v",
			resToken.OpenID, err)
		middleware.ResponseError(ctx, ecode.ServerErr, err)
		return
	}

	uri := ctx.Query("uri")
	if uri == "" {
		uri = "index.html"
	}
	url := consts.UrlPrefix + "/" + uri
	ctx.SetCookie("token", token, 7200, "", "", false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, url)

}

func (c *wechatController) Echo(ctx *gin.Context) {

	// 传入request和responseWriter
	server := c.wc.GetServer(ctx.Request, ctx.Writer)
	// 设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		// 回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	// 处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 发送回复的消息
	_ = server.Send()

}

func NewWechatController(service service.WechatService) WeChatController {
	// 创建一个wechat对象
	rdOpts := cache.RedisOpts{
		Host:        conf.C.RedisWechat.Host + ":" + strconv.Itoa(conf.C.RedisWechat.Port),
		Password:    conf.C.RedisWechat.Password,
		Database:    conf.C.RedisWechat.Db,
		MaxIdle:     conf.C.RedisWechat.MaxIdle,
		MaxActive:   conf.C.RedisWechat.MaxActive,
		IdleTimeout: int32(conf.C.RedisWechat.IdleTimeout),
	}
	redisCache := cache.NewRedis(&rdOpts)

	cfg := &wechat.Config{
		AppID:          conf.C.WeChat.AppID,
		AppSecret:      conf.C.WeChat.AppSecret,
		Token:          conf.C.WeChat.Token,
		EncodingAESKey: conf.C.WeChat.EncodingAESKey,
		PayMchID:       conf.C.WeChat.PayMchID,
		PayNotifyURL:   conf.C.WeChat.PayNotifyURL,
		PayKey:         conf.C.WeChat.PayKey,
		Cache:          redisCache,
	}

	return &wechatController{
		cfg:     cfg,
		wc:      wechat.NewWechat(cfg),
		service: service,
	}
}
