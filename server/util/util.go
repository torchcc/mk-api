package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	vali "mk-api/server/validator"
)

func IsMobile(mobile string) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(mobile)
}

// 绑定参数
func ParseRequest(c *gin.Context, input interface{}) error {
	err := c.ShouldBind(input)
	var errStr string
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errStr = vali.Translate(err.(validator.ValidationErrors))
		case *json.UnmarshalTypeError:
			unmarshalTypeError := err.(*json.UnmarshalTypeError)
			errStr = fmt.Errorf("%s 类型错误，期望类型 %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
		default:
			errStr = errors.New("unknown error").Error()
		}
		Log.Error(errStr)
		return errors.New(errStr)
	}
	return nil
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}
