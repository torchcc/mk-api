package model

import (
	"time"

	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
)

type UserAddr struct {
	Id             int64  `json:"id" db:"id"`
	UserId         int64  `json:"user_id" db:"user_id"`
	ProvinceId     int64  `json:"province_id" db:"province_id" binding:"required"`
	ProvinceName   string `json:"province_name" db:"province_name"`
	CityId         int64  `json:"city_id" db:"city_id" binding:"required"`
	CityName       string `json:"city_name" db:"city_name"`
	CountyId       int64  `json:"county_id" db:"county_id" binding:"required"`
	CountyName     string `json:"county_name" db:"county_name"`
	TownId         int64  `json:"town_id" db:"town_id" binding:"required"`
	TownName       string `json:"town_name" db:"town_name"`
	BuildingDetail string `json:"building_detail" db:"building_detail" binding:"required,min=2,max=150"`
	IsDefault      int8   `json:"is_default" db:"is_default" binding:"oneof=0 1"`
}

type UserAddrModel interface {
	FindUserAddrByUserId(userID int64) ([]UserAddr, error)
	Save(addr *UserAddr) (id int64, err error)
	CancelOriginDefaultAddr(userId int64) (err error)
	FindUserAddrByAddrId(id int64) (addr *dto.GetUserAddrOutput, err error)
	DeleteUserAddrByAddrId(id int64) (err error)
	UpdateUserAddr(id int64, addr *dto.UpdateUserAddrInput) (err error)
}

type addrDatabase struct {
	connection *sqlx.DB
}

func (db *addrDatabase) CancelOriginDefaultAddr(userId int64) (err error) {
	cmd := `UPDATE mku_user_address SET is_default = 0 WHERE user_id = ? AND is_default = 1 AND is_deleted = 0`
	_, err = db.connection.Exec(cmd, userId)
	return
}

func (db *addrDatabase) Save(addr *UserAddr) (id int64, err error) {
	cmd := `INSERT INTO mku_user_address (
					user_id, 
					province_id, 
					city_id,
					county_id,
					town_id,
					building_detail,
					create_time,
					update_time,
					is_default)
			VALUES (
					:user_id, 
					:province_id, 
					:city_id,
					:county_id,
					:town_id,
					:building_detail,
					:create_time,
					:update_time,
					:is_default)`
	rs, err := db.connection.NamedExec(cmd, map[string]interface{}{
		"user_id":         addr.UserId,
		"province_id":     addr.ProvinceId,
		"city_id":         addr.CityId,
		"county_id":       addr.CountyId,
		"town_id":         addr.TownId,
		"building_detail": addr.BuildingDetail,
		"create_time":     int(time.Now().Unix()),
		"update_time":     int(time.Now().Unix()),
		"is_default":      addr.IsDefault,
	})
	if err != nil {
		util.Log.Errorf("创建user addr  失败, err: %s", err.Error())
		return 0, err
	}
	id, err = rs.LastInsertId()
	if err != nil {
		util.Log.Errorf("创建user addr  失败, err: %s", err.Error())
		return 0, err
	}
	return id, nil
}

func (db *addrDatabase) FindUserAddrByUserId(userId int64) (addrs []UserAddr, err error) {
	// 获取用户的全部收件地址
	cmd := `SELECT 
				mua.id, 
				mua.user_id, 
				mua.province_id, 
				mua.city_id, 
				mua.county_id, 
				mua.town_id,
				mua.building_detail, 
				mua.is_default 
			FROM 
				mku_user_address AS mua 
			WHERE 
				mua.user_id = ?
				AND mua.is_deleted = 0
			ORDER BY 
				mua.is_default DESC,
				mua.update_time DESC`

	err = db.connection.Select(&addrs, cmd, userId)
	return
}

func (db *addrDatabase) FindUserAddrByAddrId(id int64) (addr *dto.GetUserAddrOutput, err error) {
	addr = new(dto.GetUserAddrOutput)
	cmd := `SELECT 
				id, 
				user_id,
				province_id,
				city_id,
				county_id,
				town_id,
				building_detail,
				is_default
			FROM
				mku_user_address
			WHERE 
				id = ? 
				AND is_deleted = 0`
	err = db.connection.Get(addr, cmd, id)
	return
}

func (db *addrDatabase) DeleteUserAddrByAddrId(id int64) (err error) {
	cmd := `UPDATE mku_user_address SET is_deleted = 1 WHERE id = ? AND is_deleted = 0`
	_, err = db.connection.Exec(cmd, id)
	return
}

func (db *addrDatabase) UpdateUserAddr(id int64, addr *dto.UpdateUserAddrInput) (err error) {
	cmd := `UPDATE 
				mku_user_address
			SET 
				province_id = :province_id,
				city_id = :city_id,
				county_id = :county_id,
				town_id = :town_id,
				building_detail = :building_detail,
				is_default = :is_default
			WHERE 
				id = :id
				AND is_deleted = 0`

	type param struct {
		Id int64 `json:"id" db:"id"`
		dto.UpdateUserAddrInput
	}
	_, err = db.connection.NamedExec(cmd, param{
		Id:                  id,
		UpdateUserAddrInput: *addr,
	})
	return
}

func NewUserAddrModel() UserAddrModel {
	return &addrDatabase{connection: dao.Db}
}
