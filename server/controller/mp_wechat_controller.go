package controller

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"
	"mk-api/server/conf"
	"mk-api/server/util"
)

// wechat 路由注册
func WeChatRegister(router *gin.RouterGroup) {
	var wechatController WeChatController = NewWechatController()
	router.GET("/", wechatController.DockWithWeChatServer)
}

type WeChatController interface {
	DockWithWeChatServer(ctx *gin.Context)
	JsApiTicket(ctx *gin.Context)
	Echo(ctx *gin.Context)
}

type wechatController struct {
}

func NewWechatController() WeChatController {
	return &wechatController{}
}

// UpdateUser godoc
// @Summary Get wechat
// @Description get a single user's info
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
	timestamp, nonce, signatureIn := ctx.Param("timestamp"), ctx.Param("nonce"), ctx.Param("signature")
	signatureGen := makeSignature(timestamp, nonce)

	if signatureGen != signatureIn {
		util.Log.Infof("signatureGen != signatureIn signatureGen=%s,signatureIn=%s\n", signatureGen, signatureIn)
		ctx.String(http.StatusOK, "%s", "wrong")

	} else {
		// 如果请求来自于微信，则原样返回echostr参数内容 以上完成后，接入验证就会生效，开发者配置提交就会成功。
		echostr := ctx.Param("echostr")
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

func (c *wechatController) JsApiTicket(ctx *gin.Context) {

}

func (c *wechatController) Echo(ctx *gin.Context) {
	// 配置微信参数
	cfg := &wechat.Config{
		AppID:          conf.C.WeChat.AppID,
		AppSecret:      conf.C.WeChat.AppSecret,
		Token:          conf.C.WeChat.Token,
		EncodingAESKey: conf.C.WeChat.EncodingAESKey,
		PayMchID:       conf.C.WeChat.PayMchID,
		PayNotifyURL:   conf.C.WeChat.PayNotifyURL,
		PayKey:         conf.C.WeChat.PayKey,
		Cache:          nil,
	}

	wc := wechat.NewWechat(cfg)
	// 传入request和responseWriter
	server := wc.GetServer(ctx.Request, ctx.Writer)
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
