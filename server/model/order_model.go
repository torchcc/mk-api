package model

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
)

type OrderModel interface {
	SaveOrder(order *dto.Order, items []*dto.OrderItem) (id int64, err error)
	UpdateOrderStatus(outTradeNo string, status int8) (err error)
	ListOrder(input *dto.ListOrderInput, userId int64) ([]*dto.ListOrderOutputEle, error)
}

type orderDatabase struct {
	connection *sqlx.DB
}

func (db *orderDatabase) ListOrder(input *dto.ListOrderInput, userId int64) ([]*dto.ListOrderOutputEle, error) {
	dict := make(map[int64]*dto.ListOrderOutputEle)
	orderIds := make([]int64, 0, 10)

	cmd1 := `
			SELECT 
				mo.id AS order_id,
				mo.out_trade_no,
				mo.status,
				mo.amount
			FROM 
			     mko_order AS mo
			WHERE 
				mo.user_id = :user_id
				AND mo.is_deleted = 0
				%s
			LIMIT :start, :offset
	`
	whereStmt := ""
	if input.Status != -1 {
		whereStmt += " AND mo.status = :status"
	}
	cmd1 = fmt.Sprintf(cmd1, whereStmt)

	start := (input.PageNo - 1) * input.PageSize
	offset := input.PageSize + 1
	params := struct {
		UserId int64 `json:"user_id" db:"user_id"`
		Start  int64 `json:"start" db:"start"`
		Offset int64 `json:"offset" db:"offset"`
		*dto.ListOrderInput
	}{
		UserId:         userId,
		Start:          start,
		Offset:         offset,
		ListOrderInput: input,
	}

	rows, err := db.connection.NamedQuery(cmd1, params)
	if err != nil {
		util.Log.Errorf("查询订单列表失败, err: [%s]", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ele dto.ListOrderOutputEle
		if err = rows.StructScan(&ele); err != nil {
			util.Log.Errorf("scan failed, err: [%s]", err)
			return nil, err
		}
		dict[ele.OrderId] = &ele
		orderIds = append(orderIds, ele.OrderId)
	}

	output := make([]*dto.ListOrderOutputEle, 0, 10)
	if len(orderIds) == 0 {
		return output, nil
	}

	aggOrderItems := make([]*dto.AggregatedOrderItem, 0, 10)
	cmd2 := `
		SELECT
			moi.pkg_id,
			moi.pkg_price,
			moi.order_id,
			mp.name AS pkg_name,
			mp.avatar_url AS pkg_avatar_url,
			moi.create_time,
			COUNT(*) AS pkg_count
		FROM
			mko_order_item AS moi
				INNER JOIN
					mkp_package AS mp
					ON moi.pkg_id = mp.id AND mp.is_deleted = 0
		WHERE
				moi.order_id IN (?)
		  AND moi.is_deleted = 0
		GROUP BY
			moi.pkg_id, moi.pkg_price, moi.create_time, mp.name, moi.order_id
	`
	cmd2, args, err := sqlx.In(cmd2, orderIds)
	if err != nil {
		return nil, err
	}
	cmd2 = db.connection.Rebind(cmd2)
	err = db.connection.Select(&aggOrderItems, cmd2, args...)
	if err != nil {
		return nil, err
	}

	for _, item := range aggOrderItems {
		orderId := item.OrderId
		if dict[orderId].AggregatedOrderItems == nil {
			dict[orderId].AggregatedOrderItems = make([]*dto.AggregatedOrderItem, 0, 4)
		}
		dict[orderId].AggregatedOrderItems = append(dict[orderId].AggregatedOrderItems, item)
	}

	for _, v := range dict {
		output = append(output, v)
	}
	return output, nil
}

func (db *orderDatabase) UpdateOrderStatus(outTradeNo string, status int8) (err error) {
	const cmd = `
			UPDATE mko_order SET 
				status = :status
			WHERE 
				out_trade_no = :out_trade_no
				AND is_deleted = 0
			`
	rs, err := db.connection.NamedExec(cmd, map[string]interface{}{
		"status":       status,
		"out_trade_no": outTradeNo,
	})
	if err != nil {
		return err
	}
	if rows, err := rs.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return errors.New("rows effected is not equal to 1")
	}
	return
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
					out_trade_no,
					user_id,
					mobile,
					open_id,
					amount,
                    remark,
					create_time,
					update_time
				) VALUES (
					:out_trade_no,
					:user_id,
					:mobile,
					:open_id,
					:amount,
				  	:remark,
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
						pkg_price,
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
						:pkg_price,
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
