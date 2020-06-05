package token

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

// 设置 user_id_token.1232: token
// 设置 token.xxx: {user_id: id, mobile: mobile}
func SetToken(token string, mobile string, userId int64, cli redis.Conn) {
	userIdTokenKey := "string.user_id_token." + strconv.FormatInt(userId, 10)

	_ = cli.Send("SETEX", userIdTokenKey, 24*3600, token)

	tokenUserInfoKey := "hash.token." + token
	_ = cli.Send("HSET", tokenUserInfoKey, "user_id", userId)
	_ = cli.Send("HSET", tokenUserInfoKey, "mobile", mobile)
	_ = cli.Send("EXPIRE", tokenUserInfoKey, 24*3600)
	_ = cli.Flush()
}

func SetOpenIdUserInfo(openIdKey string, userId int64, mobile string, cli redis.Conn) {
	_ = cli.Send("HSET", openIdKey, "user_id", userId)
	_ = cli.Send("HSET", openIdKey, "mobile", mobile)
	_ = cli.Flush()
}

func GenerateUuid() string {
	return uuid.New().String()
}
