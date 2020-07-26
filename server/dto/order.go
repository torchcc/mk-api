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
	CartId int64 `json:"cart_id"`
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
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	PrePayID  string `json:"prePayId"`
	SignType  string `json:"signType"`
	Package   string `json:"package"`
	PaySign   string `json:"paySign"`
}

type OrderItem struct {
	UserId int64 `json:"user_id" db:"user_id"`
	// 订单order_id
	OrderId int64 `json:"order_id" db:"order_id"`
	// 套餐id
	PackageId int64 `json:"pkg_id" db:"pkg_id"`
	// 套餐价格
	PackagePrice float64 `json:"pkg_price" db:"pkg_price"`
	CreateTime   int64   `json:"create_time" db:"create_time"`
	UpdateTime   int64   `json:"update_time" db:"update_time"`
	*Examinee
}

type Order struct {
	Id         int64   `json:"id" db:"id"`
	OutTradeNo string  `json:"out_trade_no" db:"out_trade_no"`
	UserId     int64   `json:"user_id" db:"user_id"`
	Mobile     string  `json:"mobile" db:"mobile"`
	OpenId     string  `json:"open_id" db:"open_id"`
	Amount     float64 `json:"amount" db:"amount"`
	Remark     string  `json:"remark" db:"remark"`
	CreateTime int64   `json:"create_time" db:"create_time"`
	UpdateTime int64   `json:"update_time" db:"update_time"`
}

type ListOrderInput struct {
	// 页码, 不传默认第一页
	PageNo int64 `json:"page_no,default=1" form:"page_no,default=1" binding:"min=1"`
	// 每页条数, 不传默认 10
	PageSize int64 `json:"page_size,default=10" form:"page_size,default=10" binding:"min=1,max=100"`
	// 订单筛选 -1-全部，0-待付款，2-待预约(指已经付款) 3-已退款 4-已关闭 5-待评价(该功能暂时disable)
	Status int8 `json:"status" db:"status" form:"status,default=-1" binding:"min=-1,max=5"`
}

type ListOrderOutputEle struct {
	// 订单id
	OrderId int64 `json:"order_id" db:"order_id"`
	// 订单号
	OutTradeNo string `json:"out_trade_no" db:"out_trade_no"`
	// 订单状态 0-待付款，2-待预约(指已经付款) 3-已退款 4-已关闭 5-待评价
	Status int8 `json:"status" db:"status"`
	// 订单总价
	Amount float64 `json:"amount" db:"amount"`
	// 订单中的套餐列表
	AggregatedOrderItems []*AggregatedOrderItem `json:"aggregated_order_items"`
}

type AggregatedOrderItem struct {
	// 套餐id
	PackageId int64 `json:"pkg_id" db:"pkg_id"`
	// 套餐所属的order_id
	OrderId int64 `json:"order_id" db:"order_id"`
	// 套餐名称
	PackageName string `json:"pkg_name" db:"pkg_name"`
	// 套餐头像
	PackageAvatarUrl string `json:"pkg_avatar_url" db:"pkg_avatar_url"`
	// 套餐数量
	PackageCount int64 `json:"pkg_count" db:"pkg_count"`
	// 套餐单价
	PackagePrice float64 `json:"pkg_price" db:"pkg_price"`
	// 创建时间
	CreateTime int64 `json:"create_time" db:"create_time"`
}

type RetrieveOrderOutput struct {
	// 订单id
	OrderId int64 `json:"order_id" db:"order_id"`
	// 订单号
	OutTradeNo string `json:"out_trade_no" db:"out_trade_no"`
	// 订单状态 0-待付款，2-待预约(指已经付款) 3-已退款 4-已关闭 5-待评价
	Status int8 `json:"status" db:"status"`
	// 订单总价
	Amount float64 `json:"amount" db:"amount"`
	// 下单人/预约人手机号
	Mobile string `json:"mobile" db:"mobile"`
	// 订单备注
	Remark string `json:"remark" db:"remark"`
	// 聚合后包括套餐项目的的order item
	AggregatedOrderItemsWithPkgItem []*AggregatedOrderItemWithPkgItem `json:"aggregated_order_items_with_pkg_item"`
}

type AggregatedOrderItemWithPkgItem struct {
	// 套餐检查项目
	PkgItems []*PkgItemName
	// 套餐检查人信息列表
	Examinees []*ExamineeInOrderItem
	// 套餐信息
	AggregatedOrderItem
}

type PkgItemName struct {
	// 属性排序， 小的排前面
	OrderNo int64 `json:"order_no" db:"order_no"`
	// 属性名字，此处是套餐项目名字
	Name string `json:"name" db:"name"`
}

type ExamineeInOrderItem struct {
	// order item(订单项)的id
	OrderItemId int64 `json:"order_item_id" db:"order_item_id"`
	Examinee
}

type OItemWithPkgBrief struct {
	OrderItemId      int64   `json:"order_item_id" db:"order_item_id"`
	PackageId        int64   `json:"pkg_id" db:"pkg_id"`
	PackagePrice     float64 `json:"pkg_price" db:"pkg_price"`
	PackageName      string  `json:"pkg_name" db:"pkg_name"`
	PackageAvatarUrl string  `json:"pkg_avatar_url" db:"pkg_avatar_url"`
	OrderId          int64   `json:"order_id" db:"order_id"`
	CreateTime       int64   `json:"create_time" db:"create_time"`
	Examinee
}

type PutOrderItemInput struct {
	// 要修改的order_item 的ID
	Id int64 `json:"order_item_id" db:"id"  binding:"required"`
	// 该order_item 所属的套餐id
	PackageId int64 `json:"pkg_id" db:"pkg_id" binding:"required"`
	// 体检人信息
	Examinee
}

type CancelOrderInput struct {
	Id int64 `json:"order_id" db:"id"  binding:"required"`
	// 取消原因id, 1-支付时出故障，支付不了， 2-付款时 余额限制了 3-买多了/不想买了 4-信息写错，重新下单 5-朋友/网上评价不好 6-计划有变，时间按排不上，7-其他
	CancelReasonId int64 `json:"cancel_reason_id" binding:"required,min=1,max=7" db:"cancel_reason_id"`
	// 取消的问题描述
	Remark string `json:"remark" db:"remark"`
}
