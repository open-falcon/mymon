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

func GlobalStatus(m *MysqlIns, db mysql.Conn) ([]*MetaData, error) {
	return mysqlState(m, db, "SHOW /*!50001 GLOBAL */ STATUS")
}

func GlobalVariables(m *MysqlIns, db mysql.Conn) ([]*MetaData, error) {
	return mysqlState(m, db, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

func mysqlState(m *MysqlIns, db mysql.Conn, sql string) ([]*MetaData, error) {
	rows, _, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	data := make([]*MetaData, len(rows))
	i := 0
	for _, row := range rows {
		key_ := row.Str(0)
		v, err := row.Int64Err(1)
		// Ignore non digital value
		if err != nil {
			continue
		}

		data[i] = NewMetric(key_)
		data[i].SetValue(v)
		i++
	}
	return data[:i], nil
}
