package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"mk-api/library/ecode"
)

type ResponseCode int

type Response struct {
	EMessage string      `json:"emsg"`
	Ecode    ecode.Code  `json:"ecode"`
	Data     interface{} `json:"data"`
}

func ResponseError(c *gin.Context, code ecode.Code, err error) {
	resp := &Response{Ecode: code, EMessage: err.Error(), Data: ""}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(200, err)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	resp := &Response{Ecode: ecode.OK, EMessage: "", Data: data}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
