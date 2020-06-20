package dto

type PostOrderInput struct {
	// 订单项目
	CartItems []*CartItem `json:"cart_items" binding:"required,dive"`
	// 预约人手机号
	SubscriberMobile  string `json:"subscriber_mobile" binding:"required,checkMobile"`
	SubscriberComment string `json:"subscriber_comment"`
}

type CartItem struct {
	// 购物车条目id
	CartId int64 `json:"cart_id" binding:"required"`
	// 套餐id
	PackageId int64 `json:"pkg_id" binding:"required"`
	// 套餐数量
	PackageCount int `json:"pkg_count" db:"pkg_count" binding:"required"`
	// 体检人信息列表
	Examinees []*Examinee `json:"examinees" binding:"omitempty,dive"`
}

type Examinee struct {
	// 体检人姓名
	ExamineeName string `json:"examinee_name" db:"examinee_name" binding:"required"`
	// 体检人电话
	ExamineeMobile string `json:"examinee_mobile" binding:"required,checkMobile" db:"examinee_mobile"`
	// 身份证号码
	IdCardNo string `json:"id_card_no" db:"id_card_no" binding:"required,checkIdCardNo"`
	// 性别 1-男 2-女
	Gender int8 `json:"gender" binding:"required" db:"gender"`
	// 婚否 1-是 2-否
	IsMarried int8 `json:"is_married" db:"is_married"`
	// 体检日期 时间戳， 精确到 日
	ExamineDate int64 `json:"examine_date" binding:"required" db:"examine_date"`
}

type PostOrderOutput struct {
	// 订单id
	OrderId int64 `json:"id"`
	// 订单流水号
	BizNo uint64 `json:"biz_no"`
	// 订单总金额 单位：分（显示到ui时需要前端转为元）
	Amount float64 `json:"amount"`
}

type OrderItem struct {
	UserId     int64 `json:"user_id" db:"user_id"`
	OrderId    int64 `json:"order_id" db:"order_id"`
	PackageId  int64 `json:"pkg_id" db:"pkg_id"`
	CreateTime int64 `json:"create_time" db:"create_time"`
	UpdateTime int64 `json:"update_time" db:"update_time"`
	*Examinee
}

type Order struct {
	Id         int64   `json:"id" db:"id"`
	BizNo      int64   `json:"biz_no" db:"biz_no"`
	UserId     int64   `json:"user_id" db:"user_id"`
	Mobile     string  `json:"mobile" db:"mobile"`
	OpenId     string  `json:"open_id" db:"open_id"`
	Amount     float64 `json:"amount" db:"amount"`
	CreateTime int64   `json:"create_time" db:"create_time"`
	UpdateTime int64   `json:"update_time" db:"update_time"`
}
