// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

var SlaveStatusToSend = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_Log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

func slaveStatus(m *MysqlIns, db mysql.Conn) ([]*MetaData, error) {

	isSlave := NewMetric("Is_slave")

	row, res, err := db.QueryFirst("SHOW SLAVE STATUS")
	if err != nil {
		return nil, err
	}

	// be master
	if row == nil {
		isSlave.SetValue(0)
		return []*MetaData{isSlave}, nil
	}

	// be slave
	isSlave.SetValue(1)

	data := make([]*MetaData, len(SlaveStatusToSend))
	for i, s := range SlaveStatusToSend {
		data[i] = NewMetric(s)
		switch s {
		case "Slave_SQL_Running", "Slave_IO_Running":
			data[i].SetValue(0)
			v := row.Str(res.Map(s))
			if v == "Yes" {
				data[i].SetValue(1)
			}
		default:
			v, err := row.Int64Err(res.Map(s))
			if err != nil {
				data[i].SetValue(-1)
			} else {
				data[i].SetValue(v)
			}
		}
	}
	return append(data, isSlave), nil
}
