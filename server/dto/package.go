package dto

type ListPackageInput struct {
	// 页码, 不传默认第一页
	PageNo int64 `json:"page_no,default=1" form:"page_no,default=1" binding:"min=1"`
	// 每页条数, 不传默认 10
	PageSize int64 `json:"page_size,default=10" form:"page_size,default=10" binding:"min=1,max=100"`

	// 医院等级， 0-不限 1-公立三甲 2-公立医院 3-民营医院 4-专业机构
	Level int8 `json:"level" form:"level" binding:"oneof=0 1 2 3 4" db:"level"`

	// 套餐分类 id, 0-不限
	CategoryId int64 `json:"category_id" form:"category_id" binding:"min=0" db:"category_id"`
	// 价格区间左值,0 表示不限 单位分
	MinPrice int64 `json:"min_price" form:"min_price" binding:"min=0,max=3000000" db:"min_price"`
	// 价格区间右值 0 表示不限 单位分
	MaxPrice int64 `json:"max_price" form:"max_price" binding:"min=0,max=100000000" db:"max_price"`
	// 适用人群 0-不限 1-男士 2-女未婚 3-女已婚
	Target int8 `json:"target" form:"target" binding:"oneof=0 1 2 3" db:"target"`
	// 检测目标高发疾病id
	DiseaseId int8 `json:"disease_id,default=0" form:"disease_id,default=0" db:"disease_id"`
	// 优先排序 0-默认排序，1-低价优先 2 高价优先
	OrderBy int8 `json:"order_by" form:"order_by" binding:"oneof=0 1 2"`
	// 按套餐名字搜索
	Name string `json:"name" form:"name" db:"name"`
}

type ListPackageOutputEle struct {
	// 套餐id
	Id         int64 `json:"id"`
	HospitalId int64 `json:"hospital_id" db:"hospital_id"`
	// 套餐名字
	Name string `json:"name"`
	// 医院体检中心名字
	HospitalName string `json:"hospital_name" db:"hospital_name"`
	// 套餐头像
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	// 医院等级， 0-不限 1-公立三甲 2-公立医院 3-民营医院 4-专业机构
	Level int8 `json:"level" binding:"oneof=0 1 2 3 4" db:"level"`
	// 已经预约的单数, 这个暂时需要前端用hidden隐藏起来
	Sold int64 `json:"sold" db:"sold"`
	// 门市价, 原价, 单位分
	PriceOriginal float64 `json:"price_original" db:"price_original"`
	// 真实价格， 现价格， 单位分
	PriceReal float64 `json:"price_real" db:"price_real"`
}

type GetPackageOutPut struct {
	// 套餐基本信息
	BasicInfo *PackageBasicInfo `json:"basic_info"`
	// 具体项目
	Items []PackageItem `json:"items"`
	// 套餐须知
	Notices []PackageNotice `json:"notices"`
	// 套餐流程
	Procedure []PackageProcedure `json:"procedure"`
}

type PackageBasicInfo struct {
	Id int64 `json:"id" db:"id"`
	// 套餐名称
	Name       string `json:"name"`
	HospitalId int64  `json:"hospital_id" db:"hospital_id"`
	// 医院体检中心名称
	HospitalName string `json:"hospital_name" db:"hospital_name"`
	// 套餐URL
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	// 原价/门市价
	PriceOriginal float64 `json:"price_original" db:"price_original"`
	// 实际价格
	PriceReal float64 `json:"price_real" db:"price_real"`
	// 已经预约的数量
	Sold int64 `json:"sold" db:"sold"`
	// 套餐简介
	Brief string `json:"brief" db:"brief"`
	// 套餐详细介绍
	Comment string `json:"comment" db:"comment"`
	// 套餐提示
	Tips string `json:"tips"`
	// 套餐目标人群 0-不限 1-男 2-未婚女 3-已婚女
	Target int8 `json:"target" db:"target"`
}

type PackageAttribute struct {
	Id    int64 `json:"id" db:"id"`
	PkgId int64 `json:"pkg_id" db:"pkg_id"`
	// 套餐的属性类型， 1-套餐项目 2-套餐须知 3-套餐体检的流程
	AttrType int8 `json:"attr_type" db:"attr_type"`
	// 属性排序， 小的排前面
	OrderNo int64 `json:"order_no" db:"order_no"`
	// 属性名字
	Name string `json:"name" db:"name"`
	// 属性概述
	Brief string `json:"brief" db:"brief"`
	// 属性详情
	Comment string `json:"comment" db:"comment"`
	// 创建时间
	CreateTime int64 `json:"create_time" db:"create_time"`
}

// 设置排序 先按照order_no 排， 再按照create_time排序, 下次试试这里是指针的情况·
type PackageAttributes []PackageAttribute

func (p PackageAttributes) Len() int {
	return len(p)
}

func (p PackageAttributes) Less(i, j int) bool {
	if p[i].OrderNo == p[j].OrderNo {
		return p[i].CreateTime < p[j].CreateTime
	}
	return p[i].OrderNo < p[j].OrderNo
}

func (p PackageAttributes) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PackageItem = PackageAttribute
type PackageNotice = PackageAttribute
type PackageProcedure = PackageAttribute

type PkgTargetNPrice struct {
	Price  float64 `db:"price_real"`
	Target int8    `db:"target"`
}

type Category struct {
	// 类别(专项疾病)id
	Id int64 `json:"id" db:"id"`
	// 类别(专项疾病)名称
	Name string `json:"name" db:"name"`
}

type Disease Category
