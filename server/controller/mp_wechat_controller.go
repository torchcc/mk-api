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

const (
	aboutAppointment = "详细请联系人工客服或致电客服热线：0668-2853837"
	aboutInvoice     = `如个人购买客户需开具发票，请点击“人工客服”后，发送开票订单号,发票抬头公司全称、税号，发票接收地址，收件人信息等信息。
发票3-5个工作日内，快递为您寄出。`
	aboutRefund = `<div class="content">
                <div style="margin-bottom: 20px;">
                  <p>改退说明</p>
                  <ul>
                    <li>
                      <span>退款：</span>
                      如客户预约成功后选择退款，需扣除套餐实付金额的10%作为服务费。
                    </li>
                    <li>
                      <span>改期:</span>
                      每位客户拥有3次改期机会，无需支付任何费用的权益。3次后再改期，需扣除套餐实付金额的10%作为服务费。
                    </li>
                    <li>
                      <span>补偿：</span>
                      客户预约好体检时间后，由迈康体检网的原因造成客户体检当天无法进行体检，迈康体检网会补偿不超过支付金额的10%作为补偿费用。
                    </li>
                  </ul>
                </div>

                <div>
                  <p>弃检变更</p>
                  <ul>
                    <li>
                      <span>弃检：</span>
                      当您预约套餐时，即表示接受检测的所有项目。如因自身原因放弃体检套餐中的检查项目，网站将不予退款处理。
                    </li>
                    <li>
                      <span>变更：</span>
                      套餐里的体检项目，可能会由于体检中心设备或其他原因，体检中心会自动帮您更换同价位其他项目或升级项目，望您能理解和支持。
                    </li>
                  </ul>
                </div>
              </div>`
	noticeBeforeExam = `<div class="content">
                <div style="margin-bottom: 20px;">
<p>体检须知</p>
<ul>
<li>1、 体检当天早晨须空腹。 建议着装休闲，宽松，方便检时操作及更换衣物。</li>
<li>2、体检前三天注意不要吃油腻、不易消化的食物体检前一天晚上8点之后不再进餐(12点之后不 可饮水),保证睡眠;避免剧烈运动和情绪激动,以保证体检结果的准确性 </li>
<li>3、参加X线检查,请勿穿着带有金属饰物或配件的衣物,孕妇及半年内准备怀孕的受检者请勿做X线检查及妇科检查。 </li>
<li>4、如检査盆腔的子宮及其附件、膀胱、前列腺等脏器时,检查前需保留膀胱尿液,可在检査前2小时饮温开水1000毫升左右,检查前2-4小时内不要小便。 </li>
<li>6、已婚女性检査妇科前需先排空尿液,经期请勿 做妇科检查,可预约时间再检查(7天之内)。未婚女性请勿做妇科检查。 </li>
<li>7、有眼压、眼底、裂隙灯检查项目请勿戴隐形眼镜,如戴隐形眼镜请自备眼药水和镜盒。</li>
                  </ul>
                </div>
              </div>`
	commonQuestion = `猜您想问
<a href="weixin://bizmsgmenu?msgmenucontent=人工&msgmenuid=999">0.人工客服</a>
<a href="weixin://bizmsgmenu?msgmenucontent=预约相关&msgmenuid=1">1.预约相关</a>
<a href="weixin://bizmsgmenu?msgmenucontent=开具发票&msgmenuid=2">2.开具发票</a>
<a href="weixin://bizmsgmenu?msgmenucontent=关于退款&msgmenuid=3">3.关于退款</a>
<a href="weixin://bizmsgmenu?msgmenucontent=检前注意事项&msgmenuid=4">4.检前注意事项</a>
`
	greetings       = "Hello， 感谢您关注迈康体检， 我们时刻关注您的健康。\n<a href='https://www.mkhealth.club/'>点次选购体检套餐和预约体检</a>\n\n如需人工服务，回复\"人工\"开始咨询，人工客服时间为7:30-20:00"
	plainExamNotice = `体检须知
1、 体检当天早晨须空腹。 建议着装休闲，宽松，方便检时操作及更换衣物。
2、体检前三天注意不要吃油腻、不易消化的食物体检前一天晚上8点之后不再进餐(12点之后不 可饮水),保证睡眠;避免剧烈运动和情绪激动,以保证体检结果的准确性 
3、参加X线检查,请勿穿着带有金属饰物或配件的衣物,孕妇及半年内准备怀孕的受检者请勿做X线检查及妇科检查。 
4、如检査盆腔的子宮及其附件、膀胱、前列腺等脏器时,检查前需保留膀胱尿液,可在检査前2小时饮温开水1000毫升左右,检查前2-4小时内不要小便。 
6、已婚女性检査妇科前需先排空尿液,经期请勿 做妇科检查,可预约时间再检查(7天之内)。未婚女性请勿做妇科检查。 
7、有眼压、眼底、裂隙灯检查项目请勿戴隐形眼镜,如戴隐形眼镜请自备眼药水和镜盒。`
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
	c.WXMsgReply(ctx, &msg)
}

func (c *wechatController) WXMsgReply(ctx *gin.Context, mixMessage *message.MixMessage) {
	if mixMessage.MsgType == message.MsgTypeEvent {
		// 关注事件发送欢迎消息
		if mixMessage.Event == message.EventSubscribe {
			replyText(ctx, mixMessage, greetings)
			return
		}

		switch mixMessage.EventKey {
		case "examineNotice":
			replyText(ctx, mixMessage, plainExamNotice)
		case "onlineCustomService":
			onlineCustomService(ctx, mixMessage)
		default:
			responseEmptyStr(ctx, mixMessage)
		}
		return
	}

	if mixMessage.MsgType == message.MsgTypeText {
		switch mixMessage.Content {
		case "人工", "客服":
			onlineCustomService(ctx, mixMessage)
		case "预约相关":
			replyText(ctx, mixMessage, aboutAppointment)
		case "开具发票":
			replyText(ctx, mixMessage, aboutInvoice)
		case "关于退款":
			replyText(ctx, mixMessage, aboutRefund)
		case "检前注意事项":
			replyText(ctx, mixMessage, noticeBeforeExam)
		default:
			replyText(ctx, mixMessage, commonQuestion) // 自助客服
		}
		return
	}
	replyText(ctx, mixMessage, commonQuestion) // 自助客服
}

// 响应空字符串,注意是响应体是空字符串，不是content是空字符串
func responseEmptyStr(ctx *gin.Context, mixMessage *message.MixMessage) {
	_, _ = ctx.Writer.Write([]byte(""))
}

func replyText(ctx *gin.Context, mixMessage *message.MixMessage, content string) {
	respMsg := message.NewText(content)
	respMsg.SetMsgType(message.MsgTypeText)
	respMsg.SetFromUserName(mixMessage.ToUserName)
	respMsg.SetToUserName(mixMessage.FromUserName)
	respMsg.SetCreateTime(time.Now().Unix())
	msg, err := xml.Marshal(respMsg)
	if err != nil {
		util.Log.Errorf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = ctx.Writer.Write(msg)
}

// 点击在线客服
func onlineCustomService(ctx *gin.Context, mixMessage *message.MixMessage) {
	tcMsg := message.NewTransferCustomer("")
	tcMsg.SetMsgType(message.MsgTypeTransfer)
	tcMsg.SetCreateTime(time.Now().Unix())
	tcMsg.SetFromUserName(mixMessage.ToUserName)
	tcMsg.SetToUserName(mixMessage.FromUserName)
	msg, err := xml.Marshal(tcMsg)
	if err != nil {
		util.Log.Errorf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = ctx.Writer.Write(msg)
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
