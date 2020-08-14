package service

import (
	"time"

	"mk-api/server/model"
)

// func startTimer(f func()) {
// 	go func() {
// 		for {
// 			now := time.Now()
// 			// 计算下一个零点
// 			next := now.Add(time.Hour * 24)
// 			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
// 			t := time.NewTimer(next.Sub(now))
// 			<-t.C
// 			f()
// 		}
// 	}()
// }

func startTimer(f func()) {
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if time.Now().Format("15") == "00" {
				f()
			}
		}
	}()
}

func init() {
	// 每天增加套餐销售量
	startTimer(model.IncreasePkgSalesVolume)
}
