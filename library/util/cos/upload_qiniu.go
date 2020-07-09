package cos

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"mk-api/library/util/conf"
)

// 接收两个参数 一个文件流 一个 bucket 你的七牛云标准空间的名字
func Upload2QiNiu(file *multipart.FileHeader) (err error, path string, key string) {
	putPolicy := storage.PutPolicy{
		Scope: conf.C.QiniuCos.Bucket,
	}
	mac := qbox.NewMac(conf.C.QiniuCos.AccessKey, conf.C.QiniuCos.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	f, e := file.Open()
	if e != nil {
		fmt.Println(e)
		return e, "", ""
	}
	dataLen := file.Size
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	err = formUploader.Put(context.Background(), &ret, upToken, fileKey, f, dataLen, &putExtra)
	if err != nil {
		fmt.Printf("upload file failed, err: [%s]", err.Error())
		return err, "", ""
	}
	return err, conf.C.QiniuCos.ImgPath + "/" + ret.Key, ret.Key
}
