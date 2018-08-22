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
	"time"

	"github.com/pkg/errors"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" //with mysql
)

// NewMySQLConnection the constructor of mysql connecting
func NewMySQLConnection(conf *Config) (mysql.Conn, error) {
	return initMySQLConnection(conf)
}

// QueryResult the result of query
func initMySQLConnection(conf *Config) (db mysql.Conn, err error) {
	db = mysql.New("tcp", "", fmt.Sprintf(
		"%s:%d", conf.DataBase.Host, conf.DataBase.Port),
		conf.DataBase.User, conf.DataBase.Password)
	db.SetTimeout(500 * time.Millisecond)
	if err = db.Connect(); err != nil {
		err = errors.Wrap(err, "Building mysql connection failed!")
	}
	return
}
