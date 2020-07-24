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
	FindPackagePriceNTargetById(id int64) (output *dto.PkgTargetNPrice, err error)
	FindPkgItemNameByPkgId(pkgId int64) ([]*dto.PkgItemName, error)
	ListDisease() ([]*dto.Disease, error)
	ListCategory() ([]*dto.Category, error)
}

type packageDatabase struct {
	connection *sqlx.DB
}

func (db *packageDatabase) ListCategory() ([]*dto.Category, error) {
	output := make([]*dto.Category, 0, 16)
	const cmd = `SELECT id, name FROM mkp_category WHERE is_deleted = 0 LIMIT 200`
	err := db.connection.Select(&output, cmd)
	return output, err
}

func (db *packageDatabase) ListDisease() ([]*dto.Disease, error) {
	output := make([]*dto.Disease, 0, 16)
	const cmd = `SELECT id, name FROM mkp_disease WHERE is_deleted = 0 LIMIT 200`
	err := db.connection.Select(&output, cmd)
	return output, err
}

func (db *packageDatabase) FindPkgItemNameByPkgId(pkgId int64) ([]*dto.PkgItemName, error) {
	names := make([]*dto.PkgItemName, 0, 16)
	const cmd = `SELECT 
					order_no, name 
				FROM
				     mkp_package_attribute 
				WHERE 
					pkg_id = ?
					AND attr_type = 1
					AND is_deleted = 0
				ORDER BY order_no
`
	err := db.connection.Select(&names, cmd, pkgId)
	return names, err
}

func (db *packageDatabase) FindPackagePriceNTargetById(id int64) (output *dto.PkgTargetNPrice, err error) {
	output = &dto.PkgTargetNPrice{}
	cmd := `SELECT price_real, target FROM mkp_package WHERE id = ? AND is_deleted = 0`
	err = db.connection.Get(output, cmd, id)
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
				mp.target,
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
	if input.Name != "" {
		whereStmt += " AND mp.name LIKE '%" + input.Name + "%' "
	}
	if input.CategoryId != 0 {
		whereStmt += " AND EXISTS(SELECT id FROM mkp_package_category mpc WHERE pkg_id=mp.id AND mpc.category_id=:category_id AND mpc.is_deleted = 0) "
	}
	if input.DiseaseId != 0 {
		whereStmt += " AND EXISTS(SELECT id FROM mkp_package_disease mpd WHERE pkg_id=mp.id AND mpd.disease_id=:disease_id AND mpd.is_deleted = 0) "
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
				mp.price_original,
				mp.price_real,
				mh.level
			FROM 
				mkp_package AS mp
			INNER JOIN 
				mkh_hospital AS mh 
					ON mp.hospital_id = mh.id AND mh.is_deleted = 0
			WHERE 
				mp.is_deleted = 0 
				%s
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
