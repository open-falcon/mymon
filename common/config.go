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
	"os"

	"github.com/go-ini/ini"
)

// BaseConf config about dir, log, etc.
type BaseConf struct {
	BaseDir      string
	SnapshotDir  string
	SnapshotDay  int
	LogDir       string
	Endpoint     string
	LogFile      string
	LogLevel     int
	FalconClient string
	IgnoreFile   string
}

// DatabaseConf config about database
type DatabaseConf struct {
	User     string
	Password string
	Host     string
	Port     int
}

// Config for initializing. This can be loaded from TOML file with -c
type Config struct {
	Base     BaseConf
	DataBase DatabaseConf
}

// NewConfig the constructor of config
func NewConfig(file string) (*Config, error) {
	conf, err := readConf(file)
	return &conf, err
}

func readConf(file string) (conf Config, err error) {
	_, err = os.Stat(file)
	if err != nil {
		file = fmt.Sprint("etc/", file)
		_, err = os.Stat(file)
		if err != nil {
			panic(err)
		}
	}
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}
	snapshotDay, err := cfg.Section("default").Key("snapshot_day").Int()
	if err != nil {
		fmt.Println("No Snapshot!")
		snapshotDay = -1
	}
	logLevel, err := cfg.Section("default").Key("log_level").Int()
	if err != nil {
		fmt.Println("Log level default: 7!")
		logLevel = 7
	}
	host := cfg.Section("mysql").Key("host").String()
	if host == "" {
		fmt.Println("Host default: 127.0.0.1!")
		host = "127.0.0.1"
	}
	port, err := cfg.Section("mysql").Key("port").Int()
	if err != nil {
		fmt.Println("Port: default 3306!")
		port = 3306
		err = nil
	}
	conf = Config{
		BaseConf{
			BaseDir:      cfg.Section("default").Key("basedir").String(),
			SnapshotDir:  cfg.Section("default").Key("snapshot_dir").String(),
			SnapshotDay:  snapshotDay,
			LogDir:       cfg.Section("default").Key("log_dir").String(),
			LogFile:      cfg.Section("default").Key("log_file").String(),
			Endpoint:     cfg.Section("default").Key("endpoint").String(),
			LogLevel:     logLevel,
			FalconClient: cfg.Section("default").Key("falcon_client").String(),
			IgnoreFile:   cfg.Section("default").Key("ignore_file").String(),
		},
		DatabaseConf{
			User:     cfg.Section("mysql").Key("user").String(),
			Password: cfg.Section("mysql").Key("password").String(),
			Host:     host,
			Port:     port,
		},
	}
	return
}
