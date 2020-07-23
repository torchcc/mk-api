package util

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"

	"mk-api/server/conf"
)

var xlog = logrus.New()
var Log *logrus.Entry

func init() {

	hooker, err := mgorus.NewHookerWithAuthDb(
		conf.C.MongoLog.Host+":"+strconv.Itoa(conf.C.MongoLog.Port),
		conf.C.MongoLog.AuthDb,
		conf.C.MongoLog.Db,
		conf.C.MongoLog.Collection,
		conf.C.MongoLog.User,
		conf.C.MongoLog.Password)

	if err == nil {
		xlog.Hooks.Add(hooker)
	} else {
		panic(err)
	}

	if BRANCH := os.Getenv("BRANCH"); BRANCH == "test" || BRANCH == "local" {
		xlog.Level = logrus.DebugLevel
	} else {
		xlog.Level = logrus.DebugLevel
	}

	// xlog.Formatter = &logrus.JSONFormatter{}

	xlog.SetReportCaller(true)

	// the default is os.Stderr
	xlog.Out = os.Stdout

	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = xlog.Out

	var hostname string
	hostname, _ = os.Hostname()

	Log = xlog.WithFields(logrus.Fields{
		"sys_name":  conf.ServiceName,
		"host_name": hostname})

}
