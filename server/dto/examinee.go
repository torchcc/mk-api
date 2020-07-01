package dto

type PostExamineeInput struct {
	// 体检人与本人的关系 0 本人， 1 父亲， 2 兄弟姐妹， 3 儿子， 4 女儿， 5 母亲， 6 夫妻， 7 其他
	Relation int8 `json:"relation,default=0" db:"relation" binding:"min=0,max=7"`
	// 体检人姓名
	ExamineeName string `json:"examinee_name" db:"examinee_name" binding:"required"`
	// 体检人电话
	ExamineeMobile string `json:"examinee_mobile" binding:"required,checkMobile" db:"examinee_mobile"`
	// 身份证号码
	IdCardNo string `json:"id_card_no" db:"id_card_no" binding:"required,checkIdCardNo"`
	// 婚否 1-是 2-否
	IsMarried int8 `json:"is_married" db:"is_married" binding:"required"`
}

type ExamineeBean struct {
	UserId int64 `json:"user_id" db:"user_id" binding:"-"`
	// 性别 1-男 2-女
	Gender     int8  `json:"gender" binding:"required" db:"gender"`
	CreateTime int64 `json:"create_time" db:"create_time" binding:"-"`
	UpdateTime int64 `json:"update_time" db:"update_time" binding:"-"`
	PostExamineeInput
}

// POST 时请忽略这个参数.

type ListExamineeOutputEle struct {
	// Examinee 表的id
	Id int64 `json:"id" db:"id"`
	// 性别 1-男 2-女
	Gender int8 `json:"gender" binding:"required" db:"gender"`
	// 年龄
	Age int64 `json:"age" binding:"-"`
	PostExamineeInput
}
