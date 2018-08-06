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
	"testing"

	"github.com/open-falcon/mymon/common"

	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	conf, err := common.NewConfig(testConfigFile)
	if err != nil {
		t.Error(err)
	}
	data := NewMetric(conf, "test")
	assert.NotNil(t, data, "Create new metric failed!")
}

func TestGetIsReadOnly(t *testing.T) {
	//test master
	isReadOnlyByMaster, _ := GetIsReadOnly(dbTestMaster)
	assert.Equal(t, isReadOnlyByMaster, 0, "Read only judge error!")
	//test slave
	isReadOnlyBySlave, _ := GetIsReadOnly(dbTestSlave)
	assert.Equal(t, isReadOnlyBySlave, 1, "Read only judge error!")
}
