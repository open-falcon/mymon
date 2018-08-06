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
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/open-falcon/mymon/common"

	"github.com/astaxie/beego/logs"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

// Global tag var
var (
	IsSlave    int
	IsReadOnly int
	Tag        string
)

//Log logger of project
var Log *logs.BeeLogger

func main() {
	// parse config file
	var confFile string
	flag.StringVar(&confFile, "c", "myMon.cfg", "myMon configure file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(fmt.Sprintf("%10s: %s", "Version", Version))
		fmt.Println(fmt.Sprintf("%10s: %s", "Compile", Compile))
		fmt.Println(fmt.Sprintf("%10s: %s", "Branch", Branch))
		fmt.Println(fmt.Sprintf("%10s: %d", "GitDirty", GitDirty))
		os.Exit(0)
	}
	conf, err := common.NewConfig(confFile)
	if err != nil {
		fmt.Printf("NewConfig Error: %s\n", err.Error())
		return
	}
	if conf.Base.LogDir != "" {
		err = os.MkdirAll(conf.Base.LogDir, 0755)
		if err != nil {
			fmt.Printf("MkdirAll Error: %s\n", err.Error())
			return
		}
	}
	if conf.Base.SnapshotDir != "" {
		err = os.MkdirAll(conf.Base.SnapshotDir, 0755)
		if err != nil {
			fmt.Printf("MkdirAll Error: %s\n", err.Error())
			return
		}
	}

	// init log and other necessary
	Log = common.MyNewLogger(conf, common.CompatibleLog(conf))

	db, err := common.NewMySQLConnection(conf)
	if err != nil {
		fmt.Printf("NewMySQLConnection Error: %s\n", err.Error())
		return
	}
	defer func() { _ = db.Close() }()

	// start...
	Log.Info("MySQL Monitor for falcon")
	go timeout()
	err = fetchData(conf, db)
	if err != nil && err != io.EOF {
		Log.Error("Error: %s", err.Error())
	}
}

func timeout() {
	time.AfterFunc(TimeOut*time.Second, func() {
		Log.Error("Execute timeout")
		os.Exit(1)
	})
}

func fetchData(conf *common.Config, db mysql.Conn) (err error) {
	defer func() {
		MySQLAlive(conf, err == nil)
	}()

	// Get GLOBAL variables
	IsReadOnly, err = GetIsReadOnly(db)
	if err != nil {
		return
	}
	Tag = GetTag(conf)

	// SHOW XXX Metric
	var data []*MetaData

	// Get slave status and set IsSlave global var
	slaveState, err := ShowSlaveStatus(conf, db)
	if err != nil {
		return
	}

	globalStatus, err := ShowGlobalStatus(conf, db)
	if err != nil {
		return
	}
	data = append(data, globalStatus...)

	globalVars, err := ShowGlobalVariables(conf, db)
	if err != nil {
		return
	}
	data = append(data, globalVars...)

	innodbState, err := ShowInnodbStatus(conf, db)
	if err != nil {
		return
	}
	data = append(data, innodbState...)

	data = append(data, slaveState...)

	binaryLogStatus, err := ShowBinaryLogs(conf, db)
	if err != nil {
		return
	}
	data = append(data, binaryLogStatus...)

	// Send Data to falcon-agent
	msg, err := SendData(conf, data)
	Log.Info("Send response %s:%d - %s", conf.DataBase.Host, conf.DataBase.Port, string(msg))

	err = ShowProcesslist(conf, db)
	if err != nil {
		return
	}
	return
}
