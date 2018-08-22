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
)

func TestLogger(t *testing.T) {
	conf, _ := NewConfig(testConfigFile)
	log := MyNewLogger(conf, "myMon.log")
	log.Error("error test ok") // level 3, show
	log.Info("info test")      // level 6, do not show
	log.Debug("%v", *conf)     // level 7, do not show
}
