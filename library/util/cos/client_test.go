package cos

import (
	"context"
	"fmt"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
	"mk-api/server/util"
)

func initUpload(c *cos.Client, name string) *cos.InitiateMultipartUploadResult {
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	util.Log.Infof("有错误: [%s]", err)
	fmt.Printf("%#v\n", v)
	return v
}

func TestClient(t *testing.T) {

	cli := NewCosClient()
	fn := "/Users/troy/Desktop/handsome.jpg"

	up := initUpload(cli, fn)
	uploadID := up.UploadID
	t.Log(uploadID)
}
