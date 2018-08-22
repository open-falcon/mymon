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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	testCase := Config{
		Base: BaseConf{
			BaseDir:      ".",
			LogDir:       "./fixtures",
			Endpoint:     "endpoint",
			LogFile:      "",
			SnapshotDir:  "./snapshot",
			IgnoreFile:   "./falconignore",
			SnapshotDay:  10,
			LogLevel:     5,
			FalconClient: "http://127.0.0.1:1988/v1/push",
		},
		DataBase: DatabaseConf{
			User:     "root",
			Password: "1tIsB1g3rt",
			Host:     "127.0.0.1",
			Port:     3306,
		},
	}
	TestFile := "../fixtures/test.cfg"
	conf, err := NewConfig(TestFile)
	assert.Nil(t, err)
	assert.Equal(t, testCase, *conf)
}
