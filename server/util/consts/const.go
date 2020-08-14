package consts

import "time"

const ServiceName = "mk-server"
const UrlPrefix = "https://www.mkhealth.club"

const (
	OrderExpireIn      = time.Second * 2 * 3600
	Closed        int8 = 4
	Success       int8 = 2
)

const (
	CacheCategory = "string.CATEGORY"
	CacheDisease  = "string.DISEASE"
	CachePackage  = "string.PACKAGE"
	CacheOrder    = "string.ORDER"
	CacheProfile  = "string.PROFILE"
)

// Api Cache Duration
const (
	CategoryListDuration = time.Minute * 15
	DiseaseListDuration  = time.Minute * 15
	PackageListDuration  = time.Minute * 10
	PackageOneDuration   = time.Minute * 10
	OrderListDuration    = time.Second * 2
	ProfileOneDuration   = time.Hour * 24
)
