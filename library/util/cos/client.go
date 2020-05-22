package cos

import (
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

const BucketUrl = "https://avatar-bucket-1302104842.cos.ap-guangzhou.myqcloud.com"

func NewCosClient() *cos.Client {
	u, _ := url.Parse(BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  C.Cos.SecretID,
			SecretKey: C.Cos.SecretKey,
		},
	})
	return client
}
