package dto

type ResourceID struct {
	Id int64 `json:"id" comment:"资源id"`
}

type PaginateListOutput struct {
	// 页码
	PageNo int64 `json:"page_no" validate:"min=1"`
	// 本页实际条目数量
	PageSize int64 `json:"page_size" validate:"min=1, max=100"`
	// 是否有下一页, 1 是， 0 否
	HasNext int8 `json:"has_next"`

	// 列表数据
	List interface{} `json:"list"`
}
