package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"mk-api/library/ecode"
	. "mk-api/server/dao"
)

// TokenAuthMiddleware 检查request header 的token， 必须是注册(绑定手机)并且登录的用户 request才能往下进行
func MobileBoundRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ResponseError(ctx, ecode.RequestErr, errors.New("缺少请求token"))
			ctx.Abort()
			return
		}

		cli := Rdb.TokenRdbP.Get()
		defer cli.Close()

		tokenUserInfoKey := "hash.token." + token
		if res, _ := cli.Do("EXISTS", tokenUserInfoKey); res.(int64) <= 0 {
			ResponseError(ctx, ecode.RequestErr, errors.New("token 已经过期失效， 请重新打开微信同意授权进入"))
			ctx.Abort()
			return
		}

		userId, _ := redis.Int64(cli.Do("HGET", tokenUserInfoKey, "user_id"))
		mobile, _ := redis.String(cli.Do("HGET", tokenUserInfoKey, "mobile"))

		if mobile == "" {
			ResponseError(ctx, ecode.MobileNoVerfiy, errors.New("用户尚未绑定手机"))
			ctx.Abort()
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func TokenRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ResponseError(ctx, ecode.RequestErr, errors.New("缺少请求token"))
			ctx.Abort()
			return
		}

		cli := Rdb.TokenRdbP.Get()
		defer cli.Close()

		tokenUserInfoKey := "hash.token." + token
		if res, _ := cli.Do("EXISTS", tokenUserInfoKey); res.(int64) <= 0 {
			ResponseError(ctx, ecode.RequestErr, errors.New("token 已经过期失效， 请重新打开微信同意授权进入"))
			ctx.Abort()
			return
		}

		userId, _ := redis.Int64(cli.Do("HGET", tokenUserInfoKey, "user_id"))
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
