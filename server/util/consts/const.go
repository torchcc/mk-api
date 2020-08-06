package consts

import "time"

const ServiceName = "mk-server"
const UrlPrefix = "https://www.mkhealth.club"

const (
	OrderExpireIn      = time.Second * 2 * 3600
	Closed        int8 = 4
	Success       int8 = 2
	Unpaid        int8 = 0
)

const (
	CacheCategory = "CATEGORY"
	CacheDisease  = "DISEASE"
	CachePackage  = "PACKAGE"
)

// Api Cache Duration
const (
	CategoryListDuration = time.Minute * 15
	DiseaseListDuration  = time.Minute * 15
	PackageListDuration  = time.Minute * 10
)
