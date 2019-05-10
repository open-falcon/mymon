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
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/open-falcon/mymon/common"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

// ShowGlobalStatus execute mysql query `SHOW GLOBAL STATUS`
func ShowGlobalStatus(conf *common.Config, db mysql.Conn) ([]*MetaData, error) {
	return parseMySQLStatus(conf, db, "SHOW /*!50001 GLOBAL */ STATUS")
}

// ShowGlobalVariables execute mysql query `SHOW GLOBAL VARIABLES`
func ShowGlobalVariables(conf *common.Config, db mysql.Conn) ([]*MetaData, error) {
	return parseMySQLStatus(conf, db, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

// ShowProcesslist execute mysql query `SHOW FULL PROCESSLIST`
func ShowProcesslist(conf *common.Config, db mysql.Conn) error {
	var note string
	rows, _, err := db.Query("SHOW FULL PROCESSLIST")
	if err != nil {
		Log.Debug("get processlist error: %+v", err)
		return err
	}
	for _, row := range rows {
		note += fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			time.Now().Format("2006-01-02 15:04:05"), row.Int(0), row.Str(1), row.Str(2),
			row.Str(3), row.Str(4), row.Int(5), row.Str(6), row.Str(7))
	}

	fileNameDay, fileNameOldDay := common.GetFileNameDayAndOldDay(conf, "process")
	err = Snapshot(conf, note, fileNameDay, fileNameOldDay)
	return err
}

// ShowInnodbStatus execute mysql query `SHOW SHOW /*!50000 ENGINE*/ INNODB STATUS`
func ShowInnodbStatus(conf *common.Config, db mysql.Conn) ([]*MetaData, error) {
	status, _, err := db.QueryFirst("SHOW /*!50000 ENGINE*/ INNODB STATUS")
	if err != nil {
		Log.Debug("show innodb status error: %+v", err)
		return nil, err
	}
	allStatus := status.Str(2)
	fileNameDay, fileNameOldDay := common.GetFileNameDayAndOldDay(conf, "innodb")
	err = Snapshot(conf, allStatus, fileNameDay, fileNameOldDay)
	if err != nil {
		Log.Debug("write snapshot error: %+v", err)
		return nil, err
	}

	rows := strings.Split(allStatus, "\n")
	return parseInnodbStatus(conf, rows)
}

// ShowBinaryLogs execute mysql query `SHOW BINARY LOGS`
func ShowBinaryLogs(conf *common.Config, db mysql.Conn) ([]*MetaData, error) {
	var sum int
	binlogFileCounts := NewMetric(conf, "binlog_file_counts")
	binlogFileSize := NewMetric(conf, "binlog_file_size")

	rows, _, err := db.Query("SHOW GLOBAL VARIABLES like 'log_bin'")
	if err != nil {
		Log.Debug("SHOW GLOBAL VARIABLES like 'log_bin' Error: %s", err.Error())
		return nil, err
	}

	log_bin_enabled := true
	for _, row := range rows {
		if row.Str(1) == "OFF" {
			log_bin_enabled = false
		}
	}

	if !log_bin_enabled {
		//log_bin is not enabled
		return nil, err
	}

	rows, res, err := db.Query("SHOW BINARY LOGS")
	if err != nil {
		Log.Debug("SHOW BINARY LOGS Error: %s", err.Error())
		return nil, err
	}

	for _, row := range rows {
		size, _ := strconv.Atoi(row.Str(res.Map("File_size")))
		sum += size
	}
	binlogFileCounts.SetValue(len(rows))
	binlogFileSize.SetValue(sum)
	return []*MetaData{binlogFileCounts, binlogFileSize}, err
}

// ShowSlaveStatus get all slave status of mysql serves
func ShowSlaveStatus(conf *common.Config, db mysql.Conn) ([]*MetaData, error) {
	// Check IsSlave
	row, res, err := db.QueryFirst("SHOW SLAVE STATUS")
	if err != nil {
		IsSlave = -1
		return nil, err
	}

	if row != nil {
		IsSlave = 1
	} else {
		IsSlave = 0
	}

	Tag = GetTag(conf)
	isSlaveMetric := NewMetric(conf, "Is_slave")
	isSlaveMetric.SetValue(IsSlave)

	// be master
	if IsSlave == 0 {
		// Master_is_readonly VS master_is_read_only for version compatible, ugly
		masterReadOnly, err := ShowOtherMetric(conf, db, "Master_is_readonly")
		if err != nil {
			Log.Debug("get Master_is_readonly metric error: %+v", err)
			return nil, err
		}
		masterReadOnly2, err := ShowOtherMetric(conf, db, "master_is_read_only")
		if err != nil {
			Log.Debug("get master_is_readonly metric error: %+v", err)
			return nil, err
		}
		innodbStatsOnMetadata, err := ShowOtherMetric(conf, db, "innodb_stats_on_metadata")
		if err != nil {
			// mysql < 5.5 doesn't have this variable
			Log.Debug("get innodb_stats_on_metadata metric error: %+v", err)
			return nil, err
		}
		return []*MetaData{isSlaveMetric, masterReadOnly, masterReadOnly2, innodbStatsOnMetadata}, nil
	}

	// be slave
	ioDelay, err := ShowOtherMetric(conf, db, "io_thread_delay")
	if err != nil {
		Log.Debug("get io_thread_delay metric error: %+v", err)
		return nil, err
	}
	slaveReadOnly, err := ShowOtherMetric(conf, db, "slave_is_read_only")
	if err != nil {
		Log.Debug("get slave_is_read_only metric error: %+v", err)
		return nil, err
	}
	heartbeat, err := ShowOtherMetric(conf, db, "Heartbeats_Behind_Master")
	if err != nil {
		// mysql.heartbeat table not necessary exist if you don't care about heartbeat
		// bypass heartbeat table not exist error
		Log.Debug("Heartbeats_Behind_Master: %s", err.Error())
		err = nil
	}
	data := make([]*MetaData, len(SlaveStatus))
	for i, s := range SlaveStatus {
		data[i] = NewMetric(conf, s)
		switch s {
		case "Slave_SQL_Running", "Slave_IO_Running":
			data[i].SetValue(0)
			v := row.Str(res.Map(s))
			if v == "Yes" {
				data[i].SetValue(1)
			}
		default:
			pos := res.Map(s)
			if pos > 0 {
				v, err := row.Int64Err(pos)
				if err != nil {
					data[i].SetValue(-1)
				} else {
					data[i].SetValue(v)
				}
			}
		}
	}
	return append(data, []*MetaData{isSlaveMetric, ioDelay, slaveReadOnly, heartbeat}...), nil
}

// ShowOtherMetric all other metric will add in this func
func ShowOtherMetric(conf *common.Config, db mysql.Conn, metric string) (*MetaData, error) {
	var err error
	var row mysql.Row
	newMetaData := NewMetric(conf, metric)
	switch metric {
	case "master_is_read_only", "slave_is_read_only", "Master_is_readonly":
		newMetaData.SetValue(IsReadOnly)
	case "innodb_stats_on_metadata":
		row, _, err = db.QueryFirst("SELECT /*!50504 @@GLOBAL.innodb_stats_on_metadata,*/ -1;")
		newMetaData.SetValue(row.Int(0))
	case "io_thread_delay":
		var res mysql.Result
		row, res, err = db.QueryFirst("SHOW SLAVE STATUS")
		if bytes.Equal([]byte(row.Str(res.Map("Master_Log_File"))), []byte(row.Str(res.Map("Relay_Master_Log_File")))) {
			newMetaData.SetValue(0)
		} else {
			masterLogFile := common.GetLastNum(row.Str(res.Map("Master_Log_File")), ".")
			relayMasterLogFile := common.GetLastNum(row.Str(res.Map("Relay_Master_Log_File")), ".")
			newMetaData.SetValue(masterLogFile - relayMasterLogFile)
			if masterLogFile < 0 || relayMasterLogFile < 0 {
				newMetaData.SetValue(-1)
			}
		}
	case "Heartbeats_Behind_Master":
		row, _, err = db.QueryFirst("select ts from mysql.heartbeat limit 1")
		// when row is empty, err is nil either
		if err != nil || len(row) == 0 {
			newMetaData.SetValue(-1)
		} else {
			localTimezone, _ := time.LoadLocation("Local")
			heartbeatTimeStr := row.Str(0)
			b := strings.Replace(heartbeatTimeStr, "T", " ", 1)
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", strings.Split(b, ".")[0], localTimezone)
			heartbeatTimestamp := t.Unix()
			currentTimestamp := time.Now().Unix()
			newMetaData.SetValue(currentTimestamp - heartbeatTimestamp)
		}
	}

	return newMetaData, err
}

func parseMySQLStatus(conf *common.Config, db mysql.Conn, sql string) ([]*MetaData, error) {
	rows, _, err := db.Query(sql)
	if err != nil {
		Log.Debug("sql query: %s error: %+v", sql, err)
		return nil, err
	}

	data := make([]*MetaData, len(rows))
	i := 0
	for _, row := range rows {
		k := row.Str(0)
		//TODO: many variables values "ON|OFF" but not int should be fetched, like super_read_only
		v, err := row.Int64Err(1)
		if err != nil {
			continue
		}

		data[i] = NewMetric(conf, k)
		data[i].SetValue(v)
		i++
	}
	return data[:i], nil
}

func parseInnodbSection(
	conf *common.Config, row string, section string,
	pdata *[]*MetaData, longTranTime *int) error {
	var err error
	switch section {
	case "TRANSACTIONS":
		txPrefixes := []string{"ACTIVE ", "ACTIVE (PREPARED) "}
		for _, txPrefix := range txPrefixes {
			if strings.Contains(row, txPrefix) {
				var tmpLongTransactionTime int
				secString := strings.Split(strings.Split(
					row, txPrefix)[1], " sec")[0]
				tmpLongTransactionTime, err = strconv.Atoi(secString)
				if err != nil {
					continue
				}
				if tmpLongTransactionTime > *longTranTime {
					*longTranTime = tmpLongTransactionTime
				}
				break
			}
		}
		if err != nil {
			Log.Warn(err.Error(), "Longest_transaction metric parse Error: ", row)
		}
		if strings.Contains(row, "History list length") {
			hisListLengthStr := strings.Split(row, "length ")[1]
			hisListLength, _ := strconv.Atoi(hisListLengthStr)
			HistoryListLength := NewMetric(conf, "History_list_length")
			HistoryListLength.SetValue(hisListLength)
			*pdata = append(*pdata, HistoryListLength)
		}
	case "SEMAPHORES":
		matches := regexp.MustCompile(`^Mutex spin waits\s+(\d+),\s+rounds\s+(\d+),\s+OS waits\s+(\d+)`).FindStringSubmatch(row)
		if len(matches) == 4 {
			spinWaits, _ := strconv.Atoi(matches[1])
			innodbMutexSpinWaits := NewMetric(conf, "Innodb_mutex_spin_waits")
			innodbMutexSpinWaits.SetValue(spinWaits)
			*pdata = append(*pdata, innodbMutexSpinWaits)

			spinRounds, _ := strconv.Atoi(matches[2])
			InnodbMutexSpinRounds := NewMetric(conf, "Innodb_mutex_spin_rounds")
			InnodbMutexSpinRounds.SetValue(spinRounds)
			*pdata = append(*pdata, InnodbMutexSpinRounds)

			osWaits, _ := strconv.Atoi(matches[3])
			InnodbMutexOsWaits := NewMetric(conf, "Innodb_mutex_os_waits")
			InnodbMutexOsWaits.SetValue(osWaits)
			*pdata = append(*pdata, InnodbMutexOsWaits)
		}
	}

	return err
}

func parseInnodbStatus(conf *common.Config, rows []string) ([]*MetaData, error) {
	var section string
	longTranTime := 0
	var err error
	var data []*MetaData
	for _, row := range rows {
		switch row {
		case "BACKGROUND THREAD":
			section = row
			continue
		case "DEAD LOCK ERRORS":
			section = row
			continue
		case "LATEST DETECTED DEADLOCK":
			section = row
			continue
		case "FOREIGN KEY CONSTRAINT ERRORS", "LATEST FOREIGN KEY ERROR":
			section = row
			continue
		case "SEMAPHORES":
			section = row
			continue
		case "TRANSACTIONS":
			section = row
			continue
		case "FILE I/O":
			section = row
			continue
		case "INSERT BUFFER AND ADAPTIVE HASH INDEX":
			section = row
			continue
		case "LOG":
			section = row
			continue
		case "BUFFER POOL AND MEMORY":
			section = row
			continue
		case "ROW OPERATIONS":
			section = row
			continue
		}
		err = parseInnodbSection(conf, row, section, &data, &longTranTime)
		if err != nil {
			Log.Debug("parse innodb section error: %+v", err)
			return nil, err
		}
	}
	longTranMetric := NewMetric(conf, "longest_transaction")
	longTranMetric.SetValue(longTranTime)
	data = append(data, longTranMetric)
	return data, nil
}
