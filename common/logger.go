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

	"github.com/astaxie/beego/logs"
)

// MyNewLogger constructor of needed logger
func MyNewLogger(conf *Config, logFile string) *logs.BeeLogger {
	return loggerInit(conf, logFile)
}

func loggerInit(conf *Config, logFile string) (log *logs.BeeLogger) {
	log = logs.NewLogger(0)
	log.EnableFuncCallDepth(true)
	log.SetLevel(conf.Base.LogLevel)
	if conf.Base.LogDir == "console" {
		_ = log.SetLogger("console")
	} else {
		_ = log.SetLogger(
			"file", fmt.Sprintf(
				`{"filename":"%s", "level":%d, "maxlines":0,
					"maxsize":0, "daily":false, "maxdays":0}`,
				logFile, conf.Base.LogLevel))
	}
	return
}
