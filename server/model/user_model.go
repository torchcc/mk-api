package model

import (
	"time"

	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/util"
)

// Data Object
type User struct {
	ID        int64  `json:"id" db:"id"`
	Mobile    string `json:"mobile"`
	UserName  string `json:"user_name" db:"user_name"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	Gender    int32  `json:"gender" db:"gender"`
	OpenId    string `json:"open_id" db:"open_id"`
	Country   string `json:"country" db:"country"`
	Province  string `json:"province" db:"province"`
	City      string `json:"city" db:"city"`
}

// Model Class
type UserModel interface {
	Save(user *User) (int64, error)
	Update(user User) error
	Delete(user User) error
	FindAll() ([]User, error)
	FindUserByID(uint32) (*User, error)
	FindUserByOpenId(openId string) (id int64, mobile string, err error)
}

type database struct {
	connection *sqlx.DB
}

func (db *database) Save(u *User) (id int64, err error) {
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

	cmd1 := `INSERT INTO mku_user (open_id, create_time, update_time) VALUES (?, ?, ?)`

	rs, err := tx.Exec(cmd1, u.OpenId, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return 0, err
	}

	id, err = rs.LastInsertId()
	if err != nil {
		return 0, err
	}

	cmd2 := `INSERT INTO mk_user_profile (
			user_id, 
			user_name, 
			avatar_url, 
			gender, 
			country, 
			province, 
			city, 
			create_time,
			update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	rs, err = tx.Exec(cmd2, id, u.UserName, u.AvatarUrl, u.Gender, u.Country, u.Province, u.City,
		time.Now().Unix(), time.Now().Unix())

	if err != nil {
		return 0, err
	}
	if _, err = rs.LastInsertId(); err != nil {
		return 0, err
	}

	return id, err
}

func (db *database) FindUserByOpenId(openId string) (id int64, mobile string, err error) {
	var u User
	cmd := `SELECT id, mobile FROM mku_user WHERE open_id = ? AND is_deleted = 0`
	err = db.connection.Get(&u, cmd, openId)
	return u.ID, u.Mobile, err
}

// model 层有错误要抛出去给 service 层
func NewUserModel() UserModel {
	return &database{
		connection: dao.Db,
	}
}

func (db *database) Update(user User) error {
	return nil

}

func (db *database) Delete(user User) error {
	return nil
}

func (db *database) FindAll() ([]User, error) {
	var users []User
	// sql
	return users, nil
}

func (db *database) FindUserByID(ID uint32) (*User, error) {
	var u User

	cmd := `SELECT
				mu.id,
				mobile,
				mup.user_name,
				avatar_url,
				gender,
				open_id  
			FROM
				mku_user AS mu          
			INNER JOIN
				mku_user_profile AS mup          
					ON mu.id = mup.user_id   
			WHERE
				mu.id = ?`
	err := db.connection.Get(&u, cmd, ID)

	return &u, err
}
