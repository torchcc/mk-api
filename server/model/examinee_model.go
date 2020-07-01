package model

import (
	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
)

type ExamineeModel interface {
	FindExamineesByUserId(userId int64) ([]*dto.ListExamineeOutputEle, error)
	SaveExaminee(examinee *dto.ExamineeBean) (id int64, err error)
	DeleteExamineeByIdNUserId(id int64, userId int64) error
	UpdateExaminee(bean *dto.ExamineeBean) error
}

type examineeDatabase struct {
	connection *sqlx.DB
}

func (db *examineeDatabase) UpdateExaminee(bean *dto.ExamineeBean) error {
	const cmd = `
		UPDATE mku_examinee SET 
			examinee_name 	  = :examinee_name   
			,relation     	  = :relation        
			,id_card_no   	  = :id_card_no      
			,is_married   	  = :is_married      
			,gender       	  = :gender          
			,examinee_mobile  = :examinee_mobile
			,update_time  = :update_time
		WHERE 
			id = :id AND user_id = :user_id AND is_deleted = 0
`
	_, err := db.connection.NamedExec(cmd, bean)
	return err
}

func (db *examineeDatabase) DeleteExamineeByIdNUserId(id int64, userId int64) error {
	const cmd = `UPDATE mku_examinee SET is_deleted = 1 WHERE id = ? AND user_id = ? and is_deleted = 0`
	_, err := db.connection.Exec(cmd, id, userId)
	return err
}

func (db *examineeDatabase) FindExamineesByUserId(userId int64) ([]*dto.ListExamineeOutputEle, error) {
	output := make([]*dto.ListExamineeOutputEle, 0, 10)
	const cmd = `SELECT
       				id
					,examinee_name  
					,relation       
					,id_card_no     
					,is_married     
					,gender         
					,examinee_mobile
				FROM mku_examinee 
				WHERE 
					user_id = ?
					AND is_deleted = 0
				ORDER BY update_time
				LIMIT 10
`
	err := db.connection.Select(&output, cmd, userId)
	return output, err
}

func (db *examineeDatabase) SaveExaminee(examinee *dto.ExamineeBean) (id int64, err error) {
	const cmd = `
			INSERT INTO mku_examinee (
				user_id        
				,examinee_name  
				,relation       
				,id_card_no     
				,is_married     
				,gender         
				,examinee_mobile
				,create_time
				,update_time
			) VALUES (
				:user_id        
				,:examinee_name  
				,:relation       
				,:id_card_no     
				,:is_married     
				,:gender         
				,:examinee_mobile
				,:create_time
				,:update_time
			)
`
	rs, err := db.connection.NamedExec(cmd, examinee)
	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}

func NewExamineeModel() ExamineeModel {
	return &examineeDatabase{connection: dao.Db}
}
