package dto

type UserDetailOutput struct {
	Id        int64  `json:"id" db:"id"`
	Mobile    string `json:"mobile"`
	UserName  string `json:"user_name" db:"user_name"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	Gender    int32  `json:"gender" db:"gender"`
	OpenId    string `json:"open_id" db:"open_id"`
}

type CreateUserAddrInput struct {
	ProvinceId     int64  `json:"province_id" db:"province_id" binding:"required" comment:"省id"`
	CityId         int64  `json:"city_id" db:"city_id" binding:"required" comment:"城市id"`
	CountyId       int64  `json:"county_id" db:"county_id" binding:"required" comment:"区id"`
	TownId         int64  `json:"town_id" db:"town_id" binding:"required" comment:"镇/街道id"`
	BuildingDetail string `json:"building_detail" db:"building_detail" binding:"required,min=2,max=150" comment:"街道楼栋门牌详细字符串"`
	IsDefault      int8   `json:"is_default" db:"is_default" binding:"oneof=0 1" comment:"是否设置为默认地址"`
}

type GetUserAddrOutput struct {
	Id             int64  `json:"id" db:"id" comment:"地址id"`
	UserId         int64  `json:"id" db:"user_id" comment:"用户id"`
	ProvinceId     int64  `json:"province_id" db:"province_id" binding:"required" comment:"省id"`
	CityId         int64  `json:"city_id" db:"city_id" binding:"required" comment:"城市id"`
	CountyId       int64  `json:"county_id" db:"county_id" binding:"required" comment:"区id"`
	TownId         int64  `json:"town_id" db:"town_id" binding:"required" comment:"镇/街道id"`
	BuildingDetail string `json:"building_detail" db:"building_detail" binding:"required,min=2,max=150" comment:"街道楼栋门牌详细字符串"`
	IsDefault      int8   `json:"is_default" db:"is_default" binding:"oneof=0 1" comment:"是否设置为默认地址"`
}

type UpdateUserAddrInput CreateUserAddrInput

type PutUserProfileInput struct {
	// 用户id
	UserId int64 `json:"user_id" db:"user_id"`
	// 用户昵称
	UserName string `json:"user_name" db:"user_name"`
	// 用户性别 1-男 2-女
	Gender int32 `json:"gender" db:"gender" binding:"oneof=0 1 2"`
	// 这个参数请忽略， 不用管
	UpdateTime int64 `json:"update_time" db:"update_time"`
}

type UploadUserAvatarOutput struct {
	// 头像url
	AvatarUrl string `json:"avatar_url"`
}

type GetUserMobileOutput struct {
	Mobile string `json:"mobile"`
}

type PostLocInput struct {
	Longitude float64 `json:"longitude" db:"longitude"`
	Latitude  float64 `json:"latitude" db:"latitude"`
}
