package model

import (
	"mk-api/server/dao"
	"mk-api/server/util"
)

/*若要在i ≤ R ≤ j 这个范围得到一个随机整数R ，需要用到表达式 FLOOR(i + RAND() * (j – i + 1))。 这里随机产生20-50*/
func IncreasePkgSalesVolume() {
	const cmd = `UPDATE mkp_package SET sold = sold + FLOOR(20 + RAND() * 31) WHERE is_deleted = 0`
	_, err := dao.Db.Exec(cmd)
	if err != nil {
		util.Log.Warning("failed to increase package sales volume!")
	} else {
		util.Log.Info("increase package sales volume done.")
	}
}
