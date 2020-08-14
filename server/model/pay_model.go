package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/silenceper/wechat/v2/pay/notify"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
)

const (
	Closed  = 4
	Success = 2
)

type PayModel interface {
	FindBillByOutTradeNo(result *notify.PaidResult) (*dto.TradeBill, error)
	SuccessPaidResult2Bill(result *notify.PaidResult) (err error)
	SaveTradeBill(bill *dto.TradeBill) (id int64, err error)
	ExpireBill(billId int64) (err error)
	CheckPayStatusByPrepayId(prepayId string) (status int8, err error)
}

type payDatabase struct {
	connection *sqlx.DB
}

func (db *payDatabase) CheckPayStatusByPrepayId(prepayId string) (status int8, err error) {
	const cmd = `SELECT status FROM mkb_trade_bill WHERE prepay_id = ? AND is_deleted = 0`
	err = db.connection.Get(&status, cmd, prepayId)
	return
}

func (db *payDatabase) ExpireBill(billId int64) (err error) {
	const cmd = `
			UPDATE mkb_trade_bill SET 
				status = :status
			WHERE 
				id = :billId
				AND time_expire < unix_timestamp(now())
				AND is_deleted = 0
`
	rs, err := db.connection.NamedExec(cmd, map[string]interface{}{
		"status": Closed,
		"billId": billId,
	})
	if err != nil {
		return err
	}
	if rows, err := rs.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.New("rows effected is not equal to 1, but " + strconv.Itoa(int(rows)))
	}
	return
}

func (db *payDatabase) SaveTradeBill(bill *dto.TradeBill) (id int64, err error) {
	const cmd = `INSERT INTO mkb_trade_bill (
					order_id      
					,out_trade_no  
					,prepay_id     
					,nonce_str
					,total_fee     
					,fee_type      
					,status        
					,trans_type    
					,time_start    
					,time_expire   
					,create_time   
					,update_time   
				) VALUES (
					:order_id      
					,:out_trade_no  
					,:prepay_id 
					,:nonce_str    
					,:total_fee     
					,:fee_type      
					,:status        
					,:trans_type    
					,:time_start    
					,:time_expire   
					,:create_time   
					,:update_time 
)
`
	rs, err := db.connection.NamedExec(cmd, bill)
	if err != nil {
		return
	}
	util.Log.Infof("success to create a trade bill, sql is: [%s], params is: [%v]", cmd, bill)
	id, err = rs.LastInsertId()
	return
}

func (db *payDatabase) SuccessPaidResult2Bill(result *notify.PaidResult) (err error) {
	const cmd = `
			UPDATE mkb_trade_bill SET 
				status = :status,
				transaction_id = :transaction_id,
				time_end = :time_end
			WHERE 
				out_trade_no = :out_trade_no
				AND is_deleted = 0
`
	timeEnd, _ := time.Parse("20060102150405", *result.TimeEnd)
	rs, err := db.connection.NamedExec(cmd, map[string]interface{}{
		"status":         Success,
		"transaction_id": *result.TransactionID,
		"time_end":       timeEnd.Unix(),
		"out_trade_no":   *result.OutTradeNo,
	})
	if err != nil {
		return err
	}
	_, err = rs.RowsAffected()
	return err
}

func (db *payDatabase) FindBillByOutTradeNo(result *notify.PaidResult) (*dto.TradeBill, error) {
	var bill dto.TradeBill
	const cmd = `
			SELECT
				id,
			    transaction_id,
			    out_trade_no,
			    time_end,
			    status,
			    total_fee
			FROM 
				mkb_trade_bill
			WHERE
				out_trade_no = ?
				AND is_deleted = 0
`
	no, err := strconv.ParseInt(*result.OutTradeNo, 10, 64)
	if err != nil {
		return nil, err
	}
	err = db.connection.Get(&bill, cmd, no)

	return &bill, err
}

func NewPayModel() PayModel {
	return &payDatabase{connection: dao.Db}
}
