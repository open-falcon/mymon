package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/open-falcon/mymon/common"
	"github.com/stretchr/testify/assert"
)

// This Test need http server to receive data
// Http URL: http://127.0.0.1:1988/v1/push
func TestFetchData(t *testing.T) {
	// init create old file to remove
	now := time.Now()
	ld := now.Add(-time.Duration(confTestMaster.Base.SnapshotDay) * 24 * time.Hour).Day()
	for _, prefix := range []string{"innodb", "process"} {
		fileName := fmt.Sprintf(
			"%s/%s_%s:%d", confTestMaster.Base.SnapshotDir, prefix,
			confTestMaster.DataBase.Host, confTestMaster.DataBase.Port)
		fileNameOldDay := fmt.Sprintf("%s_%d", fileName, ld)
		os.Create(fileNameOldDay)
	}
	// init logger
	Log = common.MyNewLogger(confTestMaster, "myMon.log")

	err := fetchData(confTestMaster, dbTestMaster)
	assert.Nil(t, err)
	assert.Equal(t, IsReadOnly, 0)
	assert.Equal(t, Tag, "port=3306,isSlave=0,readOnly=0,type=mysql")
	assert.Equal(t, IsSlave, 0)
}

// This Test need http server to receive data
// Http URL: http://127.0.0.1:1988/v1/push
func Test_Main(t *testing.T) {
	// init create old file to remove
	now := time.Now()
	ld := now.Add(-time.Duration(confTestMaster.Base.SnapshotDay) * 24 * time.Hour).Day()
	for _, prefix := range []string{"innodb", "process"} {
		fileName := fmt.Sprintf(
			"%s/%s_%s:%d", confTestMaster.Base.SnapshotDir, prefix,
			confTestMaster.DataBase.Host, confTestMaster.DataBase.Port)
		fileNameOldDay := fmt.Sprintf("%s_%d", fileName, ld)
		os.Create(fileNameOldDay)
	}
	main()
}
