/*
* Open-Falcon
*
* Copyright (c) 2014-2018 Xiaomi, Inc. All Rights Reserved.
*
* This product is licensed to you under the Apache License, Version 2.0 (the "License").
* You may not use this product except in compliance with the License.
*
* This product may include a number of subcomponents with separate copyright notices
* and license terms. Your use of these subcomponents is subject to the terms and
* conditions of the subcomponent's license, as noted in the LICENSE file.
 */
package common

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testConfigFile is used to make test
const testConfigFile = "../fixtures/test.master.cfg"

var conf *Config

func init() {
	var err error
	conf, err = NewConfig(testConfigFile)
	if err != nil {
		panic(err)
	}
}

func TestGetFilenameDayAndOldDay(t *testing.T) {
	testCase := "test"
	originStr := fmt.Sprintf(
		"%s/%s_%s:%d_", conf.Base.SnapshotDir, testCase,
		conf.DataBase.Host, conf.DataBase.Port)
	fileNameDay, fnameOldDay := GetFileNameDayAndOldDay(conf, "test")
	d, _ := strconv.Atoi(strings.Split(fileNameDay, originStr)[1])
	ld, _ := strconv.Atoi(strings.Split(fnameOldDay, originStr)[1])
	if d > ld {
		assert.Equal(t, conf.Base.SnapshotDay, d-ld)
	}
}

func TestGetLastNum(t *testing.T) {
	testString := "tests$22"
	testString2 := "tests$test"
	assert.Equal(t, 22, GetLastNum(testString, "$"))
	assert.Equal(t, -1, GetLastNum(testString, "1"))
	assert.Equal(t, -2, GetLastNum(testString2, "$"))
}

func TestHostname(t *testing.T) {
	assert.NotNil(t, Hostname(conf))
}

func TestCompatibleLog(t *testing.T) {
	testCases := []*Config{
		{Base: BaseConf{LogDir: "."}},
		{Base: BaseConf{LogDir: ".", LogFile: "./test_log"}},
		{Base: BaseConf{LogFile: "./test_log"}},
		{},
	}
	expectRes := []string{
		"./myMon.log", "./test_log", "./test_log", "myMon.log",
	}
	for i, logConf := range testCases {
		assert.Equal(t, expectRes[i], CompatibleLog(logConf))
	}
}
