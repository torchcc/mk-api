package model

import (
	"github.com/jmoiron/sqlx"
	"mk-api/server/dao"
	"mk-api/server/util"
	"sync"
)

type RegionIdName struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db "name"`
}

type RegionModel interface {
	GetRegionIdNameMap() (id2name map[int64]string, err error)
}

type regionDatabase struct {
	connection *sqlx.DB
}

var regionId2NameMap map[int64]string
var once sync.Once

// check lock check的once 实现单例模式
func (db regionDatabase) GetRegionIdNameMap() (id2name map[int64]string, err error) {
	once.Do(func() {
		var idNames []RegionIdName
		cmd := `SELECT id, name FROM mkm_region WHERE is_deleted = 0`
		err = db.connection.Select(&idNames, cmd)
		if err != nil {
			util.Log.Errorf("查询region出错，err: [%s]", err)
			return
		}
		regionId2NameMap = make(map[int64]string)
		for _, idName := range idNames {
			regionId2NameMap[idName.Id] = idName.Name
		}
	})
	return regionId2NameMap, err
}

func NewRegionModel() RegionModel {
	return &regionDatabase{connection: dao.Db}
}
