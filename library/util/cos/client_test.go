package cos

import (
	"os"
	"testing"
)

// 上传的时候指定 存储的文件名
func TestClient(t *testing.T) {

	// cli := NewCosClient("https://common-1302104842.cos.ap-guangzhou.myqcloud.com")
	//
	// // 1. 上传本地文件。
	// fileSource := "./mm.jpeg"
	// _, err := cli.Object.PutFromFile(context.Background(), "mm.jpeg", fileSource, nil)
	// if err != nil {
	// 	t.Errorf("failed to upload local image: %v\n", err)
	// }
	//
	// // 2 上传io流
	// fileSource = "./gg.jpeg"
	// f, _ := os.Open(fileSource)
	// defer f.Close()
	//
	// _, err = cli.Object.Put(context.Background(), "gg.png", f, nil)
	// if err != nil {
	// 	t.Errorf("failed to upload image in byte format : %v\n", err)
	// }
	t.Log("if it is needed to test, put an image on the cur dir as the fileSource")
}

func TestUploadIOStream(t *testing.T) {
	// fileSource := "./mm.jpeg"
	f, _ := os.Open(fileSource)
	// defer f.Close()
	// name, urlPrefix, err := UploadIOStream(fileSource, f, true)
	// if err != nil {
	// 	t.Errorf("failed to upload image using UploadIOStream func: %s\n", err.Error())
	// } else {
	// 	t.Logf("uploaded! the url of the image is: %s/%s\n", urlPrefix, name)
	// }
	// return
	t.Log("if it is needed to test, put an image on the cur dir as the fileSource")
}
