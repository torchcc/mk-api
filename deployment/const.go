package deployment

import (
	"errors"
	"os"
	"path"
	"runtime"
)

const superconfJsonFn = "superconf.json"

var (
	DEPLOY_DIR              = ""
	SUPERCONF_JSON_ABS_PATH = ""
	BRANCH                  = ""
)

func init() {
	getBranchName()
	getDeployDir()
	SUPERCONF_JSON_ABS_PATH = path.Join(DEPLOY_DIR, BRANCH, superconfJsonFn)
}

func getBranchName() {
	BRANCH = os.Getenv("BRANCH")
	if BRANCH == "" {
		BRANCH = "test"
	}
}

func getDeployDir() {
	DEPLOY_DIR = path.Dir(getCurrentFile())

	if DEPLOY_DIR == "" {
		panic(errors.New("can not get current file info"))
	}
}

func getCurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("can not get current file info"))
	}
	return file
}
