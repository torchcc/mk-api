package cos

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/tencentyun/cos-go-sdk-v5"
)

const CommonBucketUrl = "https://common-1302104842.cos.ap-guangzhou.myqcloud.com"

func NewCosClient(bucketUrl string) *cos.Client {
	u, _ := url.Parse(bucketUrl)
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

// 若不想对保存的文件取hash命名， hashName取false
func UploadIOStream(name string, r io.ReadSeeker, hashName bool) (fn string, urlPrefix string, err error) {
	cli := NewCosClient(CommonBucketUrl)

	if hashName {
		h := md5.New()
		if _, err = io.Copy(h, r); err != nil {
			return name, "", err
		}
		name = fmt.Sprintf("%x%s", h.Sum(nil), path.Ext(name))
		_, _ = r.Seek(0, io.SeekStart)
	}

	_, err = cli.Object.Put(context.Background(), name, r, nil)
	return name, CommonBucketUrl, err
}
