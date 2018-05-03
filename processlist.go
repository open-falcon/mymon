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
	"strings"
	"github.com/ziutek/mymysql/mysql"

	_ "github.com/ziutek/mymysql/native"
	log "github.com/Sirupsen/logrus"
	"reflect"
)

// ProcsStatusToSend 统计 show processlist 状态
var ProcsStatusToSend = map[string]int{
	"State_closing_tables":       0,
	"State_copying_to_tmp_table": 0,
	"State_end":                  0,
	"State_freeing_items":        0,
	"State_init":                 0,
	"State_locked":               0,
	"State_login":                0,
	"State_preparing":            0,
	"State_reading_from_net":     0,
	"State_sending_data":         0,
	"State_sorting_result":       0,
	"State_statistics":           0,
	"State_updating":             0,
	"State_writing_to_net":       0,
	"State_none":                 0,
	"State_other":                0,
}

func procsStatus(m *MysqlIns, db mysql.Conn) ([]*MetaData, error) {

	rows, _, err := db.Query("SHOW PROCESSLIST")
	if err != nil {
		return nil, err
	}

	log.Debugf("show processlist size: %d", len(rows))

	for _, row := range rows {
		key_ := row.Str(6)

		log.Debugf("key: %s, type: %s", key_, reflect.TypeOf(key_))
		if key_ == "" {
			key_ = "none"
		}

		if match("^Table lock|Waiting for .*lock$", key_) {
			key_ = "Locked"
		}

		key_ = strings.ToLower(key_)
		key_ = strings.Replace(key_, " ", "_", -1)

		if _, ok := ProcsStatusToSend["State_"+key_]; ok {
			ProcsStatusToSend["State_"+key_]++
		} else {
			ProcsStatusToSend["State_other"]++
		}

	}

	log.Debugf("ProcsStatusToSend: %v", ProcsStatusToSend)

	data := make([]*MetaData, len(ProcsStatusToSend))
	i := 0

	for k, v := range ProcsStatusToSend {
		data[i] = NewMetric(k)
		data[i].SetValue(v)
		i++
	}

	return data[:i], nil
}
