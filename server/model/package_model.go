package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
)

type PackageModel interface {
	ListPackage(input *dto.ListPackageInput) ([]dto.ListPackageOutputEle, error)
	FindPackageBasicInfo(id int64) (*dto.PackageBasicInfo, error)
	FindPackageAttr(pkgId int64) (attrs []dto.PackageAttribute, err error)
	FindPackagePriceById(id int64) (price float64, err error)
}

type packageDatabase struct {
	connection *sqlx.DB
}

func (db *packageDatabase) FindPackagePriceById(id int64) (price float64, err error) {
	cmd := `SELECT price_real FROM mkp_package WHERE id = ? AND is_deleted = 0`
	err = db.connection.Get(&price, cmd, id)
	return
}

func (db *packageDatabase) FindPackageAttr(pkgId int64) ([]dto.PackageAttribute, error) {
	var attrs = make([]dto.PackageAttribute, 0, 0)
	cmd := `SELECT
				id,
				pkg_id,
				attr_type,
				order_no,
				name, 
				brief,
				comment,
				create_time
			FROM 
				mkp_package_attribute 
			WHERE
				pkg_id = ?
				AND is_deleted = 0
			`
	err := db.connection.Select(&attrs, cmd, pkgId)
	return attrs, err

}

func (db *packageDatabase) FindPackageBasicInfo(id int64) (*dto.PackageBasicInfo, error) {
	var basicInfo dto.PackageBasicInfo
	cmd := `SELECT
				mp.id, 
				mp.hospital_id,
				mp.name,
				mh.name AS hospital_name,
				mp.avatar_url,
				mp.price_original,
				mp.price_real,
				mp.brief,
				mp.comment,
				mp.tips
			FROM 
				mkp_package AS mp 
			INNER JOIN 
				mkh_hospital AS mh
					ON mp.hospital_id = mh.id
					AND mh.is_deleted = 0
			WHERE 
				mp.id = ?
				AND mp.is_deleted = 0
			`
	err := db.connection.Get(&basicInfo, cmd, id)
	return &basicInfo, err
}

func (db *packageDatabase) ListPackage(input *dto.ListPackageInput) ([]dto.ListPackageOutputEle, error) {
	elems := make([]dto.ListPackageOutputEle, 0, 0)

	start := (input.PageNo - 1) * input.PageSize
	offset := input.PageSize + 1
	whereStmt := ""
	orderByStmt := ""
	if input.Level != 0 {
		whereStmt += " AND mh.level = :level "
	}
	if input.CategoryId != 0 {
		whereStmt += " AND mc.id = :category_id"
	}
	if input.MinPrice != 0 {
		whereStmt += " AND mp.price_real >= :min_price "
	}
	if input.MaxPrice != 0 {
		whereStmt += " AND mp.price_real <= :max_price "
	}
	if input.Target != 0 {
		whereStmt += " AND mp.target = :target"
	}
	switch input.OrderBy {
	case 1:
		orderByStmt = " ORDER BY mp.price_real "
	case 2:
		orderByStmt = " ORDER BY mp.price_real DESC "
	}

	cmd := ` SELECT 
				mp.id,
				mp.hospital_id,
				mp.name,
				mh.name AS hospital_name,
				mp.avatar_url,
				mc.name AS category_name,
				mc.id AS category_id,
				mh.level
			FROM 
				mkp_package AS mp
			INNER JOIN 
				mkh_hospital AS mh 
					ON mp.hospital_id = mh.id AND mh.is_deleted = 0
			INNER JOIN 
				mkp_package_category AS mpc 
					ON mp.id = mpc.pkg_id AND mpc.is_deleted = 0 
			INNER JOIN 
				mkp_category AS mc 
					ON mc.id = mpc.category_id AND mc.is_deleted = 0
			WHERE 
				mp.is_deleted = 0 
				%s
			-- GROUP BY mc.name
				%s
			LIMIT :start, :offset
`
	cmd = fmt.Sprintf(cmd, whereStmt, orderByStmt)
	util.Log.Infof("查询套餐列表sql: [%s], 参数: [%v]", cmd, input)
	params := struct {
		Start  int64 `json:"start"`
		Offset int64 `json:"offset"`
		dto.ListPackageInput
	}{
		Start: start, Offset: offset, ListPackageInput: *input,
	}
	rows, err := db.connection.NamedQuery(cmd, params)
	if err != nil {
		util.Log.Errorf("查询套餐列表出错, sql: [%s], 参数： [%v], err: [%s]", cmd, input, err.Error())
		return elems, err
	}
	defer rows.Close()
	for rows.Next() {
		var p dto.ListPackageOutputEle
		err = rows.StructScan(&p)
		if err != nil {
			util.Log.Errorf("scan failed, err: [%s]", err.Error())
			return elems, err
		}
		elems = append(elems, p)
	}
	return elems, err

}

func NewPackageModel() PackageModel {
	return &packageDatabase{connection: dao.Db}
}
