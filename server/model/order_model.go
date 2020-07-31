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
	FindOrderDetailById(id int64, pkgModel PackageModel) (*dto.RetrieveOrderOutput, error)
	DeleteOrderByIdNUserId(userId int64, id int64) error
	FindOrderPayStatusById(orderId int64) (*dto.OrderPayStatus, error)
	UpdateOrderItem(input *dto.PutOrderItemInput) error
	CancelOrder(input *dto.CancelOrderInput) error
	RefundOrder(input *dto.RefundOrderInput) (int64, error)
	FindOutTradeNoByOrderId(orderId int64) (string, error)
	FindRefundReasonIdByOrderId(orderId int64) int64
	FindOrderInfo2NotifyClientByOutTradeNo(outTradeNo string) *dto.OInfo4PaidNotify
}

type orderDatabase struct {
	connection *sqlx.DB
}

func (db *orderDatabase) FindOrderInfo2NotifyClientByOutTradeNo(outTradeNo string) *dto.OInfo4PaidNotify {
	var output *dto.OInfo4PaidNotify
	const cmd = `SELECT id, open_id, out_trade_no, amount FROM mko_order WHERE out_trade_no = ? AND is_deleted = 0`
	_ = db.connection.Get(output, cmd, outTradeNo)
	return output
}

func (db *orderDatabase) FindRefundReasonIdByOrderId(orderId int64) int64 {
	var refundReasonId int64
	const cmd = `SELECT refund_reason_id FROM mko_order WHERE id = ? AND is_deleted = 0`
	_ = db.connection.Get(&refundReasonId, cmd, orderId)
	return refundReasonId
}

func (db *orderDatabase) FindOutTradeNoByOrderId(orderId int64) (string, error) {
	var outTradeNo string
	const cmd = `SELECT out_trade_no FROM mko_order WHERE id = ? AND is_deleted = 0`
	err := db.connection.Get(&outTradeNo, cmd, orderId)
	return outTradeNo, err
}

func (db *orderDatabase) RefundOrder(input *dto.RefundOrderInput) (int64, error) {
	const cmd = `
			UPDATE mko_order SET 
				refund_reason_id = :refund_reason_id,
				refund_reason_remark = :refund_reason_remark
			WHERE id = :id AND is_deleted = 0 AND status = 2
`
	rs, err := db.connection.NamedExec(cmd, input)
	if err != nil {
		errStr := fmt.Sprintf("failed to update refund reason, err: [%s]", err)
		return 0, errors.New(errStr)
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		errStr := fmt.Sprintf("failed to get rows effected, err: [%s]", err)
		return 0, errors.New(errStr)
	}
	return rows, err

}

func (db *orderDatabase) CancelOrder(input *dto.CancelOrderInput) error {
	cmd := `
			UPDATE mko_order SET
			status = 4,
			cancel_reason_id = :cancel_reason_id,
			%s
			update_time = UNIX_TIMESTAMP(NOW())
			WHERE 
			id = :id AND status = 0 AND is_deleted = 0
`
	remarkStmt := ""
	if input.Remark != "" {
		remarkStmt = "remark = :remark, "
	}
	cmd = fmt.Sprintf(cmd, remarkStmt)
	_, err := db.connection.NamedExec(cmd, input)
	return err
}

func (db *orderDatabase) UpdateOrderItem(input *dto.PutOrderItemInput) error {
	const cmd = `
				UPDATE mko_order_item SET 
					examinee_name = :examinee_name,
					examinee_mobile = :examinee_mobile,
					id_card_no = :id_card_no,
					gender = :gender,
					is_married = :is_married,
					examine_date = :examine_date,
					update_time = UNIX_TIMESTAMP(NOW())
				WHERE id = :id AND is_deleted = 0
`
	_, err := db.connection.NamedExec(cmd, input)
	return err
}

func (db *orderDatabase) FindOrderPayStatusById(orderId int64) (*dto.OrderPayStatus, error) {
	var output dto.OrderPayStatus
	const cmd = `SELECT 
					mo.status,
					mb.prepay_id,
					mb.nonce_str,
					mb.time_expire
				FROM 
					mko_order AS mo 
					INNER JOIN mkb_trade_bill AS mb ON mo.id = mb.order_id 
				WHERE 
					mo.id = ?
					AND mo.is_deleted = 0
					AND mb.is_deleted = 0
`
	err := db.connection.Get(&output, cmd, orderId)
	return &output, err
}

func (db *orderDatabase) DeleteOrderByIdNUserId(userId int64, id int64) error {
	const cmd = `UPDATE mko_order SET is_deleted = 1 WHERE id = ? AND user_id = ?`
	_, err := db.connection.Exec(cmd, id, userId)
	return err
}

func (db *orderDatabase) FindOrderDetailById(id int64, pkgModel PackageModel) (*dto.RetrieveOrderOutput, error) {
	output := dto.RetrieveOrderOutput{}
	// step 1 获取订单表头信息
	const cmd1 = `
			SELECT
				mo.id AS order_id,
				mo.out_trade_no,
				mo.status,
				mo.amount,
				mo.mobile,
				mo.remark
			FROM 
				mko_order AS mo
			WHERE 
				mo.id = ? 
				AND mo.is_deleted = 0
`
	if err := db.connection.Get(&output, cmd1, id); err != nil {
		return nil, err
	}
	output.AggregatedOrderItemsWithPkgItem = make([]*dto.AggregatedOrderItemWithPkgItem, 0, 4)

	var orderItems []*dto.OItemWithPkgBrief
	const cmd2 = `
			SELECT 
				moi.id AS order_item_id,
				moi.pkg_id,
				moi.pkg_price,
				moi.order_id,
				moi.examinee_name,
				moi.examinee_mobile,
				moi.id_card_no,
				moi.is_married,
				moi.gender,
				moi.examine_date,
				mp.name AS pkg_name,
				mp.avatar_url AS pkg_avatar_url,
				moi.create_time
			FROM
				mko_order_item AS moi
				INNER JOIN
					mkp_package AS mp
					ON moi.pkg_id = mp.id AND mp.is_deleted = 0
			WHERE 
				moi.order_id = ?
				AND moi.is_deleted = 0
`
	if err := db.connection.Select(&orderItems, cmd2, id); err != nil {
		return nil, err
	}
	if orderItems == nil {
		return &output, nil
	}

	dic := make(map[int64]*dto.AggregatedOrderItemWithPkgItem)
	for _, item := range orderItems {
		if _, ok := dic[item.PackageId]; !ok {
			pkgItems, err := pkgModel.FindPkgItemNameByPkgId(item.PackageId)
			if err != nil {
				util.Log.Errorf("查询套餐项目名称失败, err: [%s]", err.Error())
			}
			dic[item.PackageId] = &dto.AggregatedOrderItemWithPkgItem{
				PkgItems: pkgItems,
				Examinees: []*dto.ExamineeInOrderItem{{
					OrderItemId: item.OrderItemId,
					Examinee: dto.Examinee{
						ExamineeName:   item.ExamineeName,
						ExamineeMobile: item.ExamineeMobile,
						IdCardNo:       item.IdCardNo,
						Gender:         item.Gender,
						IsMarried:      item.IsMarried,
						ExamineDate:    item.ExamineDate,
					},
				}},
				AggregatedOrderItem: dto.AggregatedOrderItem{
					PackageId:        item.PackageId,
					OrderId:          item.OrderId,
					PackageName:      item.PackageName,
					PackageAvatarUrl: item.PackageAvatarUrl,
					PackageCount:     1,
					PackagePrice:     item.PackagePrice,
					CreateTime:       item.CreateTime,
				},
			}
		} else {
			dic[item.PackageId].PackageCount++
			dic[item.PackageId].Examinees = append(dic[item.PackageId].Examinees, &dto.ExamineeInOrderItem{
				OrderItemId: item.OrderItemId,
				Examinee: dto.Examinee{
					ExamineeName:   item.ExamineeName,
					ExamineeMobile: item.ExamineeMobile,
					IdCardNo:       item.IdCardNo,
					Gender:         item.Gender,
					IsMarried:      item.IsMarried,
					ExamineDate:    item.ExamineDate,
				},
			})
		}
	}

	for _, item := range dic {
		output.AggregatedOrderItemsWithPkgItem = append(output.AggregatedOrderItemsWithPkgItem, item)
	}

	return &output, nil

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
			util.Log.Panicf("save order 的时候panic了， %#v", p)
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			util.Log.Errorf("failed to save order rolling back, err is %s", err.Error())
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
