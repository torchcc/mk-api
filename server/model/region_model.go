package model

import (
	"strconv"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"mk-api/server/dao"
	"mk-api/server/dto"
	"mk-api/server/util"
)

type RegionIdName struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type RegionModel interface {
	GetRegionIdNameMap() (id2name map[int64]string, err error)
	FindRegionsByParentId(parentId int64) ([]*dto.Region, error)
}

type regionDatabase struct {
	connection *sqlx.DB
	goCache    *cache.Cache
}

func (db regionDatabase) FindRegionsByParentId(parentId int64) (output []*dto.Region, err error) {
	// 增加内存缓存
	key := "region" + strconv.FormatInt(parentId, 10)
	if x, found := db.goCache.Get(key); found {
		output = x.([]*dto.Region)
		return
	}

	const cmd = `SELECT id, name, parent_id, level FROM mkm_region WHERE parent_id = ? AND is_deleted = 0`
	err = db.connection.Select(&output, cmd, parentId)

	if err != nil {
		util.Log.Errorf("failed to get regions, sql: %s, parent_id: %d, err: %s", cmd, parentId, err.Error())
	} else {
		db.goCache.Set(key, output, cache.DefaultExpiration)
	}
	return
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
	return &regionDatabase{
		connection: dao.Db,
		goCache:    dao.GoCache,
	}
}
