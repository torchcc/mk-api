package model

import (
	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/dto"
)

type ExamineeModel interface {
	FindExamineesByUserId(userId int64) ([]*dto.ListExamineeOutputEle, error)
	SaveExaminee(examinee *dto.ExamineeBean) (id int64, err error)
}

type examineeDatabase struct {
	connection *sqlx.DB
}

// TODO 到这儿了
func (db *examineeDatabase) FindExamineesByUserId(userId int64) ([]*dto.ListExamineeOutputEle, error) {
	output := make([]*dto.ListExamineeOutputEle, 0, 10)
	const cmd = `SELECT 
					examinee_name  
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
