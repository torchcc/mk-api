package dao

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewGoCache() *cache.Cache {
	return cache.New(15*time.Minute, 20*time.Minute)
}
