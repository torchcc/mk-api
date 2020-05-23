package sms

import (
	"encoding/json"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
	. "mk-api/library/util/conf"
)

const (
	ENDPOINT    = "sms.tencentcloudapi.com"
	SMS_SIGN    = "迈康体检网"
	msgDuration = "3" //
)

type smsParam struct {
	// 格式 ["+8618520456660", ..., ]
	PhoneNumberSet   []string
	TemplateID       string   // zk 读取
	Sign             string   // 固定为"迈康体检网"
	TemplateParamSet []string // 模板参数
	SmsSdkAppid      string   // zk读取

}

func SendRegisterMsg(mobile string, smsVerificationCode string) (err error) {

	credential := common.NewCredential(
		C.Cos.SecretID,
		C.Cos.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = ENDPOINT
	client, _ := sms.NewClient(credential, C.Cos.Region, cpf)

	request := sms.NewSendSmsRequest()

	payload := smsParam{
		PhoneNumberSet:   []string{"+86" + mobile},
		TemplateID:       C.RegisterSmsMsgTemplate.TemplateID,
		Sign:             SMS_SIGN,
		TemplateParamSet: []string{smsVerificationCode, msgDuration},
		SmsSdkAppid:      C.RegisterSmsMsgTemplate.SmsSdkAppid,
	}

	params, _ := json.Marshal(payload)
	err = request.FromJsonString(string(params))
	if err != nil {
		return
	}
	_, err = client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return
	}
	if err != nil {
		return
	}
	return
}
