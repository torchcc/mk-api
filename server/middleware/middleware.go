package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mk-api/library/ecode"
	vx "mk-api/server/validator"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "OPTIONS" {
			ctx.Next()
		} else {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			ctx.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			ctx.Header("Content-Type", "application/json")
			ctx.AbortWithStatus(200)
		}
		ctx.Next()
	}

}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		if ctx.Request.TLS != nil {
			ctx.Header("Strict-Transport-Security", "max-age=31536000")
		}
		ctx.Next()
		// Also consider adding Content-Security-Policy headers
		// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
	}
}

// Handle Errors
func HandleErrors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errorToPrint := ctx.Errors.ByType(gin.ErrorTypePublic).Last()
		if errorToPrint != nil {
			if errs, ok := errorToPrint.Err.(validator.ValidationErrors); ok {
				// trans,_ := h.uni.GetTranslator("zh") // 这里也可以通过获取 HTTP Header 中的 Accept-Language 来获取用户的语言设置

				resp := &Response{Ecode: ecode.RequestErr,
					EMessage: fmt.Sprintf("%v", errs.Translate(vx.Trans)), Data: ""}
				ctx.JSON(200, resp)
				return
			}
			// deal with other errors ...
		}
	}
}
