package dto

type GetCartOutputElem struct {
	// 购物车条目id
	Id int64 `json:"id" db:"id"`
	// 套餐头像url
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	// 套餐id
	PackageId int64 `json:"pkg_id" db:"pkg_id"`
	// 套餐名字
	PackageName string `json:"name" db:"pkg_name"`
	// 医院id
	HospitalId int64 `json:"hospital_id" db:"hospital_id"`
	// 医院名字
	HospitalName string `json:"hospital_name" db:"hospital_name"`
	// 套餐数量
	PackageCount int64 `json:"pkg_count" db:"pkg_count"`
	// 套餐单价
	PackagePrice float64 `json:"pkg_price" db:"pkg_price"`
	// 更新时间
	UpdateTime int64 `json:"update_time" db:"update_time"`
}

type PostCartInput struct {
	// 要添加的套餐id
	PkgId int64 `json:"pkg_id"`
}

type DeleteCartEntriesInput struct {
	CartIds []int64 `json:"cart_ids" db:"cart_ids"`
}
