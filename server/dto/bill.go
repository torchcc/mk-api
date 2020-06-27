package dto

type TradeBill struct {
	Id            int64  `json:"id" db:"id"`
	TransactionId string `json:"transaction_id" db:"transaction_id"`
	OrderId       int64  `json:"order_id" db:"order_id"`
	OutTradeNo    int64  `json:"out_trade_no" db:"out_trade_no"`
	PrepayId      string `json:"prepay_id" db:"prepay_id"`
	NonceStr      string `json:"nonce_str" db:"nonce_str"`
	TotalFee      int64  `json:"total_fee" db:"total_fee"`
	FeeType       int8   `json:"fee_type" db:"fee_type"`
	Status        int8   `json:"status" db:"status"`
	TransType     int8   `json:"trans_type" db:"trans_type"`
	TimeStart     int64  `json:"time_start" db:"time_start"`
	TimeExpire    int64  `json:"time_expire" db:"time_expire"`
	TimeEnd       int64  `json:"time_end" db:"time_end"`
	CreateTime    int64  `json:"create_time" db:"create_time"`
	UpdateTime    int64  `json:"update_time" db:"update_time"`
}

type UpdateBillInput struct {
	OutTradeNo    int64  `json:"out_trade_no" db:"out_trade_no"`
	TransactionId string `json:"transaction_id" db:"transaction_id"`
	TimeEnd       int64  `json:"time_end" db:"time_end"`
	UpdateTime    int64  `json:"update_time" db:"update_time"`
	Status        int8   `json:"status" db:"status"`
}

type CheckPayStatusOutput struct {
	// 支付状态 0-待支付 2-支付成功 4-订单已关闭
	Status int8 `json:"status"`
}
