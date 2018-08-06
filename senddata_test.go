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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterIgnoreData(t *testing.T) {
	testDataMetrics := []*MetaData{
		//no ignore
		{
			Metric:      "testMetric",
			Endpoint:    "localhost",
			CounterType: "type",
			Value:       "test",
			Tags:        "tag",
			Timestamp:   66,
			Step:        60,
		},
		//ignore: longest_transaction/port=3306$100
		{
			Metric:      "longest_transaction",
			Endpoint:    "localhost",
			CounterType: "type",
			Value:       1,
			Tags:        "port=3306",
			Timestamp:   66,
			Step:        60,
		},
		//ignore: Com_insert_select/isSlave=0,port=3306,readOnly=0,type=mysql
		{
			Metric:      "Com_insert_select",
			Endpoint:    "localhost",
			CounterType: "type",
			Value:       1,
			Tags:        "isSlave=0,port=3306,type=mysql,readOnly=0",
			Timestamp:   66,
			Step:        60,
		},
	}
	metricsExpectResult := []*MetaData{
		{
			Metric:      "testMetric",
			Endpoint:    "localhost",
			CounterType: "type",
			Value:       "test",
			Tags:        "tag",
			Timestamp:   66,
			Step:        60,
		},
		{
			Metric:      "longest_transaction",
			Endpoint:    "localhost",
			CounterType: "type",
			//change value
			Value:     "100",
			Tags:      "port=3306",
			Timestamp: 66,
			Step:      60,
		},
		{
			//change name
			Metric:      "_Com_insert_select",
			Endpoint:    "localhost",
			CounterType: "type",
			Value:       1,
			Tags:        "isSlave=0,port=3306,type=mysql,readOnly=0",
			Timestamp:   66,
			Step:        60,
		},
	}
	resDataMetrics := filterIgnoreData(confTestMaster, testDataMetrics)
	assert.EqualValues(t, metricsExpectResult, resDataMetrics)
}
func TestTagSame(t *testing.T) {
	tests := []struct {
		x      string
		y      string
		result bool
	}{
		{x: "a,b,c", y: "c,b,a", result: true},
		{x: "a,b", y: "a,b,c", result: false},
		{x: "a,b, c", y: "a,b,d", result: false},
	}

	for _, test := range tests {
		assert.Equal(t, test.result, tagSame(test.x, test.y), "Wrong compare tag")
	}
}

func TestSnapshot(t *testing.T) {
	testNote := "test note string!"
	testOldDayFile := "fixtures/old_day_file"
	f, err := os.Create(testOldDayFile)
	assert.Nil(t, err)
	defer f.Close()
	testDayFile := "fixtures/day_file"

	err = Snapshot(confTestMaster, testNote, testDayFile, testOldDayFile)
	assert.Nil(t, err)
	resNote, err := ioutil.ReadFile(testDayFile)
	assert.Nil(t, err)
	assert.Equal(t, testNote, string(resNote))
	err = os.Remove(testDayFile)
	assert.Nil(t, err)
}
