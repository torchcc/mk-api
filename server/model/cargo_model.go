package model

import (
	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
)

type CargoModel interface {
	FindCargoByUserId(userId int64) ([]dto.GetCargoOutputElem, error)
	IncrementPkgCount(id int64) (err error)
	FindCargoItemId(userId int64, pkgId int64) (id int64)
	CreateCargo(userId int64, pkgId int64) (err error)
	RemoveCargoEntries(input *dto.DeleteCargoEntriesInput) error
}

type cargoDatabase struct {
	connection *sqlx.DB
}

func (db *cargoDatabase) RemoveCargoEntries(input *dto.DeleteCargoEntriesInput) error {
	cmd, args, err := sqlx.In(`UPDATE mko_cargo SET is_deleted = 1 WHERE id IN (?)`, input.CargoIds)
	cmd = db.connection.Rebind(cmd)
	_, err = db.connection.Exec(cmd, args...)
	return err
}

func (db *cargoDatabase) CreateCargo(userId int64, pkgId int64) (err error) {
	cmd := `INSERT INTO mko_cargo 
				(user_id, pkg_id, pkg_count, create_time, update_time)
			VALUES 
				(?, ?, 1, UNIX_TIMESTAMP(NOW()), UNIX_TIMESTAMP(NOW()))`
	_, err = db.connection.Exec(cmd, userId, pkgId)
	return
}

func (db *cargoDatabase) IncrementPkgCount(id int64) (err error) {
	cmd := `UPDATE mko_cargo SET pkg_count = pkg_count + 1 WHERE id = ? AND is_deleted = 0`
	_, err = db.connection.Exec(cmd, id)
	return
}

func (db *cargoDatabase) FindCargoItemId(userId int64, pkgId int64) (id int64) {
	cmd := `SELECT id FROM mko_cargo WHERE user_id = ? AND pkg_id = ? AND is_deleted = 0`
	if err := db.connection.Get(&id, cmd, userId, pkgId); err != nil {
		return 0
	}
	return
}

func (db *cargoDatabase) FindCargoByUserId(userId int64) ([]dto.GetCargoOutputElem, error) {
	var pkgs = make([]dto.GetCargoOutputElem, 0, 0)
	cmd := `SELECT 
				mc.id,
				pkg_id,
				pkg_count,
				mp.price_real AS pkg_price,
				mp.avatar_url,
				mp.name AS pkg_name,
				mh.id AS hospital_id,
				mh.name AS hospital_name,
				mc.update_time
			FROM 
				mko_cargo AS mc
			INNER JOIN
				mkp_package AS mp
					ON mc.pkg_id = mp.id AND mp.is_deleted = 0
			INNER JOIN 
				mkh_hospital AS mh
					ON mp.hospital_id = mh.id AND mh.is_deleted = 0
			WHERE
				mc.user_id = ?
				AND mc.is_deleted = 0
			`
	err := db.connection.Select(&pkgs, cmd, userId)
	return pkgs, err
}

func NewCargoModel() CargoModel {
	return &cargoDatabase{connection: dao.Db}
}
