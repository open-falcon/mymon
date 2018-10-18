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

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/open-falcon/mymon/common"

	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
	"github.com/ziutek/mymysql/mysql"
)

// testConfigFile is used to make test
const testConfigFile = "fixtures/test.master.cfg"
const testSlaveConfigFile = "fixtures/test.slave.cfg"

var confTestMaster, confTestSlave *common.Config
var dbTestMaster, dbTestSlave mysql.Conn

func init() {
	var err error
	confTestMaster, err = common.NewConfig(testConfigFile)
	if err != nil {
		panic(err)
	}
	dbTestMaster, err = common.NewMySQLConnection(confTestMaster)
	if err != nil {
		panic(err)
	}

	confTestSlave, err = common.NewConfig(testSlaveConfigFile)
	if err != nil {
		panic(err)
	}
	dbTestSlave, err = common.NewMySQLConnection(confTestSlave)
	if err != nil {
		panic(err)
	}
}

/**************** TEST SHOW GLOBAL *******************/
func TestShowGlobalStatus(t *testing.T) {
	data, err := ShowGlobalStatus(confTestMaster, dbTestMaster)
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
}

func TestShowGlobalVariables(t *testing.T) {
	data, err := ShowGlobalStatus(confTestMaster, dbTestMaster)
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
}

/**************** TEST INNODB*******************/
func TestParseInnodbStatus(t *testing.T) {
	inputbuf, err := ioutil.ReadFile("fixtures/test_innodb_source.golden")
	assert.Nil(t, err)
	rows := strings.Split(string(inputbuf), "\n")

	//Test parseInnodbStatus
	res, err := parseInnodbStatus(confTestMaster, rows)
	assert.Nil(t, err)
	var actualRes string
	for i, m := range res {
		fmtLine := fmt.Sprintf("%d: MetaData Metric:%s Endpoint:%s Value:%v CounterType:%s Tags:%s Timestamp:%d Step:%d\n",
			i, m.Metric, "localhost", m.Value, m.CounterType, "test_tag", 0, m.Step)
		actualRes += fmtLine
	}
	goldie.Assert(t, "test_innodb_result", []byte(actualRes))
}

/**************** TEST BINARY *******************/
func TestShowBinaryLogStatus(t *testing.T) {
	data, err := ShowBinaryLogs(confTestMaster, dbTestMaster)
	assert.Nil(t, err)
	for _, eachData := range data {
		assert.NotNilf(t, eachData, "binlog value cannot be %v", eachData.Value)
	}
}

func TestShowSlaveStatus(t *testing.T) {
	var res []*MetaData
	var err error
	tempIsSlave := IsSlave
	defer func() { IsSlave = tempIsSlave }()

	//test master
	res, err = ShowSlaveStatus(confTestMaster, dbTestMaster)
	assert.Equal(t, 0, IsSlave)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res), "Wrong length of Metric with master!")

	//test slave
	res, err = ShowSlaveStatus(confTestSlave, dbTestSlave)
	assert.Equal(t, 1, IsSlave)
	assert.NoError(t, err)
	assert.Equal(t, len(SlaveStatus)+4, len(res), "Wrong length of Metric with slave!")
}

func TestShowOtherMetric(t *testing.T) {
	tempIsReadOnly := IsReadOnly
	defer func() { IsReadOnly = tempIsReadOnly }()
	testMetaDataStr := []string{
		"master_is_read_only", "Master_is_readonly", "slave_is_read_only",
		"innodb_stats_on_metadata", "io_thread_dela", "Heartbeats_Behind_Master", "others"}
	expectIntValues := []int{0, 0, 1}

	// test master readonly
	IsReadOnly = 0
	for i, meta := range testMetaDataStr[:2] {
		getData, _ := ShowOtherMetric(confTestMaster, dbTestMaster, meta)
		assert.Equal(t, expectIntValues[i], getData.Value, "Wrong value!")
	}
	// test slave readonly
	IsReadOnly = 1
	getData, _ := ShowOtherMetric(confTestSlave, dbTestSlave, testMetaDataStr[2])
	assert.Equal(t, expectIntValues[2], getData.Value, "Wrong value!")

	for i, meta := range testMetaDataStr[3:6] {
		getData, _ = ShowOtherMetric(confTestMaster, dbTestMaster, meta)
		assert.NotNil(t, expectIntValues[i], getData.Value, "Wrong value!")
	}
	getData, _ = ShowOtherMetric(confTestMaster, dbTestMaster, testMetaDataStr[6])
	assert.Nil(t, getData.Value, "Wrong default value of metric!")
}
