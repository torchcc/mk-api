package model

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
)

type OrderModel interface {
	SaveOrder(order *dto.Order, items []*dto.OrderItem) (id int64, err error)
}

type orderDatabase struct {
	connection *sqlx.DB
}

func (db *orderDatabase) SaveOrder(order *dto.Order, items []*dto.OrderItem) (id int64, err error) {
	tx, err := db.connection.Beginx()
	if err != nil {
		util.Log.Errorf("begin trans failed, err: %v", err)
		return 0, err
	}

	defer func() { // shi
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			util.Log.Info("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			util.Log.Info("commit")
		}
	}()

	const cmd1 = `INSERT INTO mko_order (
					biz_no,
					user_id,
					mobile,
					open_id,
					amount,
					create_time,
					update_time
				) VALUES (
					:biz_no,
					:user_id,
					:mobile,
					:open_id,
					:amount,
					:create_time,
					:update_time
				)`
	rs, err := tx.NamedExec(cmd1, order)
	if err != nil {
		return 0, err
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		item.OrderId = id
	}

	const cmd2 = `INSERT INTO mko_order_item (
						user_id,
						order_id,
						pkg_id,
						examinee_name,
						examinee_mobile,
						id_card_no,
						is_married,
						gender,
						examine_date, 
						create_time,
						update_time
					) VALUES (
						:user_id,
						:order_id,
						:pkg_id,
						:examinee_name,
						:examinee_mobile,
						:id_card_no,
						:is_married,
						:gender,
						:examine_date, 
						:create_time,
						:update_time
						)
	`
	rs, err = tx.NamedExec(cmd2, items)

	if err != nil {
		return
	}
	if rows, err := rs.RowsAffected(); err != nil {
		return 0, err
	} else if int(rows) != len(items) {
		return 0, errors.New("exec cmd2 failed")
	}

	return id, err
}

func NewOrderModel() OrderModel {
	return &orderDatabase{connection: dao.Db}
}
