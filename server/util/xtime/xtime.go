package xtime

import (
	"time"
)

func TomorrowStartAt() int64 {

	// year :=time.Now().Format("2006")
	// month := time.Now().Format("01")
	// day:=time.Now().Day()+1
	// tm2, _ := time.Parse("01/02/2006", month+"/"+strconv.Itoa(day)+"/"+year)
	// //这个算的是 距离凌晨还有多久
	// // dayLongTime:=tm2.Unix()-(3600*8)-time.Now().Unix()
	//
	// return tm2.Unix()

	timeStr := time.Now().Format("2006-01-02")

	// 使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	return t.Unix() + 1
}
