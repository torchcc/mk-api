package dto

type Region struct {
	// 省，市， 区/县， 镇/街道 的行政区域id
	Id int64 `json:"id" db:"id"`
	// 行政区域名称
	Name string `json:"name" db:"name"`
	// 该行政区域的父级id， 1级行政区域的父id 是 0
	ParentId int64 `json:"parent_id" db:"parent_id"`
	// 行政区域等级, 1-省，2-市， 3-区/县， 4-镇/街道
	Level int8 `json:"level" db:"level"`
}
