package model

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
	tokenUtil "mk-api/server/util/token"
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
	FindUserByID(id int64) (*dto.UserDetailOutput, error)
	FindUserByOpenId(openId string) (id int64, mobile string, err error)
	AddRegisterInfo(input *dto.LoginRegisterInput, userId int64) (err error)
	GetOpenIdByUserId(userId int64) (openId string, err error)
	UpdateRedisToken(openId string, userId int64, mobile string) (token string)
}

type userDatabase struct {
	connection *sqlx.DB
	redisPool  *redis.Pool
}

func (db *userDatabase) UpdateRedisToken(openId string, userId int64, mobile string) (token string) {
	cli := db.redisPool.Get()
	defer cli.Close()

	openIdKey := "hash.open_id." + openId
	tokenUtil.SetOpenIdUserInfo(openIdKey, userId, mobile, cli)
	token = tokenUtil.GenerateUuid()
	tokenUtil.SetToken(token, mobile, userId, cli)
	return
}

func (db *userDatabase) GetOpenIdByUserId(userId int64) (openId string, err error) {
	cmd := `SELECT open_id FROM mku_user WHERE id = ? AND is_deleted = 0`
	err = db.connection.Get(&openId, cmd, userId)
	return
}

func (db *userDatabase) AddRegisterInfo(input *dto.LoginRegisterInput, userId int64) (err error) {
	cmd1 := `UPDATE 
				mku_user
			SET 
				mobile = ?,
				update_time = ?
			WHERE 
				id = ?
			AND is_deleted = 0`
	_, err = db.connection.Exec(cmd1, input.Mobile, time.Now().Unix(), userId)
	if err != nil {
		return
	}

	cmd2 := `UPDATE
				mku_user_profile
			SET 
				longitude = ?,
				latitude = ?,
				update_time = ?
			WHERE 
				user_id = ?
				AND is_deleted = 0`
	_, _ = db.connection.Exec(cmd2, input.Longitude, input.Latitude, time.Now().Unix(), userId)
	return
}

func (db *userDatabase) Save(u *User) (id int64, err error) {
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

func (db *userDatabase) FindUserByOpenId(openId string) (id int64, mobile string, err error) {
	var u User
	cmd := `SELECT id, mobile FROM mku_user WHERE open_id = ? AND is_deleted = 0`
	err = db.connection.Get(&u, cmd, openId)
	return u.ID, u.Mobile, err
}

func (db *userDatabase) Update(user User) error {
	return nil

}

func (db *userDatabase) Delete(user User) error {
	return nil
}

func (db *userDatabase) FindAll() ([]User, error) {
	var users []User
	// sql
	return users, nil
}

func (db *userDatabase) FindUserByID(ID int64) (*dto.UserDetailOutput, error) {
	var u dto.UserDetailOutput

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

// model 层有错误要抛出去给 service 层
func NewUserModel() UserModel {
	return &userDatabase{
		connection: dao.Db,
		redisPool:  dao.Rdb.TokenRdbP,
	}
}
