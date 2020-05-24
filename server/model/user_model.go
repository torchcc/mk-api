package model

import (
	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
)

// Data Object
type User struct {
	ID        uint32 `json:"id"`
	Mobile    string `json:"mobile"`
	UserName  string `json:"user_name" db:"user_name"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	Gender    string `json:"gender"`
	Address   string `json:"address"`
	OpenId    string `json:"open_id" db:"open_id"`
}

// Model Class
type UserModel interface {
	Save(user User) (uint32, error)
	Update(user User) error
	Delete(user User) error
	FindAll() ([]User, error)
	FindUserByID(uint32) (*User, error)
}

type database struct {
	connection *sqlx.DB
}

// model 层有错误要抛出去给 service 层
func NewUserModel() UserModel {
	return &database{
		connection: dao.Db,
	}
}

func (db *database) Save(user User) (uint32, error) {
	// 此处写原生sql， 调用db.connection.Exec*(),
	return 1, nil
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
