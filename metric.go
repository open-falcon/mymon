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
	"fmt"
	"io"
	"time"

	"github.com/open-falcon/mymon/common"

	"github.com/ziutek/mymysql/mysql"
)

//DataType all variables should be monitor
var DataType = map[string]string{
	"Innodb_buffer_pool_reads":          DeltaPs,
	"Innodb_buffer_pool_read_requests":  DeltaPs,
	"Innodb_buffer_pool_write_requests": DeltaPs,
	"Innodb_compress_time":              DeltaPs,
	"Innodb_data_fsyncs":                DeltaPs,
	"Innodb_data_read":                  DeltaPs,
	"Innodb_data_reads":                 DeltaPs,
	"Innodb_data_writes":                DeltaPs,
	"Innodb_data_written":               DeltaPs,
	"Innodb_last_checkpoint_at":         DeltaPs,
	"Innodb_log_flushed_up_to":          DeltaPs,
	"Innodb_log_sequence_number":        DeltaPs,
	"Innodb_mutex_os_waits":             DeltaPs,
	"Innodb_mutex_spin_rounds":          DeltaPs,
	"Innodb_mutex_spin_waits":           DeltaPs,
	"Innodb_pages_flushed_up_to":        DeltaPs,
	"Innodb_rows_deleted":               DeltaPs,
	"Innodb_rows_inserted":              DeltaPs,
	"Innodb_rows_locked":                DeltaPs,
	"Innodb_rows_modified":              DeltaPs,
	"Innodb_rows_read":                  DeltaPs,
	"Innodb_rows_updated":               DeltaPs,
	"Innodb_row_lock_time":              DeltaPs,
	"Innodb_row_lock_waits":             DeltaPs,
	"Innodb_uncompress_time":            DeltaPs,

	"Binlog_event_count": DeltaPs,
	"Binlog_number":      DeltaPs,
	"Slave_count":        DeltaPs,

	"Com_admin_commands":             DeltaPs,
	"Com_assign_to_keycache":         DeltaPs,
	"Com_alter_db":                   DeltaPs,
	"Com_alter_db_upgrade":           DeltaPs,
	"Com_alter_event":                DeltaPs,
	"Com_alter_function":             DeltaPs,
	"Com_alter_procedure":            DeltaPs,
	"Com_alter_server":               DeltaPs,
	"Com_alter_table":                DeltaPs,
	"Com_alter_tablespace":           DeltaPs,
	"Com_analyze":                    DeltaPs,
	"Com_begin":                      DeltaPs,
	"Com_binlog":                     DeltaPs,
	"Com_call_procedure":             DeltaPs,
	"Com_change_db":                  DeltaPs,
	"Com_change_master":              DeltaPs,
	"Com_check":                      DeltaPs,
	"Com_checksum":                   DeltaPs,
	"Com_commit":                     DeltaPs,
	"Com_create_db":                  DeltaPs,
	"Com_create_event":               DeltaPs,
	"Com_create_function":            DeltaPs,
	"Com_create_index":               DeltaPs,
	"Com_create_procedure":           DeltaPs,
	"Com_create_server":              DeltaPs,
	"Com_create_table":               DeltaPs,
	"Com_create_trigger":             DeltaPs,
	"Com_create_udf":                 DeltaPs,
	"Com_create_user":                DeltaPs,
	"Com_create_view":                DeltaPs,
	"Com_dealloc_sql":                DeltaPs,
	"Com_delete":                     DeltaPs,
	"Com_delete_multi":               DeltaPs,
	"Com_do":                         DeltaPs,
	"Com_drop_db":                    DeltaPs,
	"Com_drop_event":                 DeltaPs,
	"Com_drop_function":              DeltaPs,
	"Com_drop_index":                 DeltaPs,
	"Com_drop_procedure":             DeltaPs,
	"Com_drop_server":                DeltaPs,
	"Com_drop_table":                 DeltaPs,
	"Com_drop_trigger":               DeltaPs,
	"Com_drop_user":                  DeltaPs,
	"Com_drop_view":                  DeltaPs,
	"Com_empty_query":                DeltaPs,
	"Com_execute_sql":                DeltaPs,
	"Com_flush":                      DeltaPs,
	"Com_grant":                      DeltaPs,
	"Com_ha_close":                   DeltaPs,
	"Com_ha_open":                    DeltaPs,
	"Com_ha_read":                    DeltaPs,
	"Com_help":                       DeltaPs,
	"Com_insert":                     DeltaPs,
	"Com_insert_select":              DeltaPs,
	"Com_install_plugin":             DeltaPs,
	"Com_kill":                       DeltaPs,
	"Com_load":                       DeltaPs,
	"Com_lock_tables":                DeltaPs,
	"Com_optimize":                   DeltaPs,
	"Com_preload_keys":               DeltaPs,
	"Com_prepare_sql":                DeltaPs,
	"Com_purge":                      DeltaPs,
	"Com_purge_before_date":          DeltaPs,
	"Com_release_savepoint":          DeltaPs,
	"Com_rename_table":               DeltaPs,
	"Com_rename_user":                DeltaPs,
	"Com_repair":                     DeltaPs,
	"Com_replace":                    DeltaPs,
	"Com_replace_select":             DeltaPs,
	"Com_reset":                      DeltaPs,
	"Com_resignal":                   DeltaPs,
	"Com_revoke":                     DeltaPs,
	"Com_revoke_all":                 DeltaPs,
	"Com_rollback":                   DeltaPs,
	"Com_rollback_to_savepoint":      DeltaPs,
	"Com_savepoint":                  DeltaPs,
	"Com_select":                     DeltaPs,
	"Com_set_option":                 DeltaPs,
	"Com_signal":                     DeltaPs,
	"Com_show_authors":               DeltaPs,
	"Com_show_binlog_events":         DeltaPs,
	"Com_show_binlogs":               DeltaPs,
	"Com_show_charsets":              DeltaPs,
	"Com_show_collations":            DeltaPs,
	"Com_show_contributors":          DeltaPs,
	"Com_show_create_db":             DeltaPs,
	"Com_show_create_event":          DeltaPs,
	"Com_show_create_func":           DeltaPs,
	"Com_show_create_proc":           DeltaPs,
	"Com_show_create_table":          DeltaPs,
	"Com_show_create_trigger":        DeltaPs,
	"Com_show_databases":             DeltaPs,
	"Com_show_engine_logs":           DeltaPs,
	"Com_show_engine_mutex":          DeltaPs,
	"Com_show_engine_status":         DeltaPs,
	"Com_show_events":                DeltaPs,
	"Com_show_errors":                DeltaPs,
	"Com_show_fields":                DeltaPs,
	"Com_show_function_status":       DeltaPs,
	"Com_show_grants":                DeltaPs,
	"Com_show_keys":                  DeltaPs,
	"Com_show_master_status":         DeltaPs,
	"Com_show_open_tables":           DeltaPs,
	"Com_show_plugins":               DeltaPs,
	"Com_show_privileges":            DeltaPs,
	"Com_show_procedure_status":      DeltaPs,
	"Com_show_processlist":           DeltaPs,
	"Com_show_profile":               DeltaPs,
	"Com_show_profiles":              DeltaPs,
	"Com_show_relaylog_events":       DeltaPs,
	"Com_show_slave_hosts":           DeltaPs,
	"Com_show_slave_status":          DeltaPs,
	"Com_show_status":                DeltaPs,
	"Com_show_storage_engines":       DeltaPs,
	"Com_show_table_status":          DeltaPs,
	"Com_show_tables":                DeltaPs,
	"Com_show_triggers":              DeltaPs,
	"Com_show_variables":             DeltaPs,
	"Com_show_warnings":              DeltaPs,
	"Com_slave_start":                DeltaPs,
	"Com_slave_stop":                 DeltaPs,
	"Com_stmt_close":                 DeltaPs,
	"Com_stmt_execute":               DeltaPs,
	"Com_stmt_fetch":                 DeltaPs,
	"Com_stmt_prepare":               DeltaPs,
	"Com_stmt_reprepare":             DeltaPs,
	"Com_stmt_reset":                 DeltaPs,
	"Com_stmt_send_long_data":        DeltaPs,
	"Com_truncate":                   DeltaPs,
	"Com_uninstall_plugin":           DeltaPs,
	"Com_unlock_tables":              DeltaPs,
	"Com_update":                     DeltaPs,
	"Com_update_multi":               DeltaPs,
	"Com_xa_commit":                  DeltaPs,
	"Com_xa_end":                     DeltaPs,
	"Com_xa_prepare":                 DeltaPs,
	"Com_xa_recover":                 DeltaPs,
	"Com_xa_rollback":                DeltaPs,
	"Com_xa_start":                   DeltaPs,
	"Com_alter_user":                 DeltaPs,
	"Com_get_diagnostics":            DeltaPs,
	"Com_lock_tables_for_backup":     DeltaPs,
	"Com_lock_binlog_for_backup":     DeltaPs,
	"Com_purge_archived":             DeltaPs,
	"Com_purge_archived_before_date": DeltaPs,
	"Com_show_client_statistics":     DeltaPs,
	"Com_show_function_code":         DeltaPs,
	"Com_show_index_statistics":      DeltaPs,
	"Com_show_procedure_code":        DeltaPs,
	"Com_show_slave_status_nolock":   DeltaPs,
	"Com_show_table_statistics":      DeltaPs,
	"Com_show_thread_statistics":     DeltaPs,
	"Com_show_user_statistics":       DeltaPs,
	"Com_unlock_binlog":              DeltaPs,

	"Aborted_clients":                    DeltaPs,
	"Aborted_connects":                   DeltaPs,
	"Access_denied_errors":               DeltaPs,
	"Binlog_bytes_written":               DeltaPs,
	"Binlog_cache_disk_use":              DeltaPs,
	"Binlog_cache_use":                   DeltaPs,
	"Binlog_stmt_cache_disk_use":         DeltaPs,
	"Binlog_stmt_cache_use":              DeltaPs,
	"Bytes_received":                     DeltaPs,
	"Bytes_sent":                         DeltaPs,
	"Connections":                        DeltaPs,
	"Created_tmp_disk_tables":            DeltaPs,
	"Created_tmp_files":                  DeltaPs,
	"Created_tmp_tables":                 DeltaPs,
	"Handler_delete":                     DeltaPs,
	"Handler_read_first":                 DeltaPs,
	"Handler_read_key":                   DeltaPs,
	"Handler_read_last":                  DeltaPs,
	"Handler_read_next":                  DeltaPs,
	"Handler_read_prev":                  DeltaPs,
	"Handler_read_rnd":                   DeltaPs,
	"Handler_read_rnd_next":              DeltaPs,
	"Handler_update":                     DeltaPs,
	"Handler_write":                      DeltaPs,
	"Opened_files":                       DeltaPs,
	"Opened_tables":                      DeltaPs,
	"Opened_table_definitions":           DeltaPs,
	"Qcache_hits":                        DeltaPs,
	"Qcache_inserts":                     DeltaPs,
	"Qcache_lowmem_prunes":               DeltaPs,
	"Qcache_not_cached":                  DeltaPs,
	"Queries":                            DeltaPs,
	"Questions":                          DeltaPs,
	"Select_full_join":                   DeltaPs,
	"Select_full_range_join":             DeltaPs,
	"Select_range_check":                 DeltaPs,
	"Select_scan":                        DeltaPs,
	"Slow_queries":                       DeltaPs,
	"Sort_merge_passes":                  DeltaPs,
	"Sort_range":                         DeltaPs,
	"Sort_rows":                          DeltaPs,
	"Sort_scan":                          DeltaPs,
	"Table_locks_immediate":              DeltaPs,
	"Table_locks_waited":                 DeltaPs,
	"Threads_created":                    DeltaPs,
	"Rpl_semi_sync_master_net_wait_time": DeltaPs,
	"Rpl_semi_sync_master_net_waits":     DeltaPs,
	"Rpl_semi_sync_master_no_times":      DeltaPs,
	"Rpl_semi_sync_master_no_tx":         DeltaPs,
	"Rpl_semi_sync_master_yes_tx":        DeltaPs,
	"Rpl_semi_sync_master_tx_wait_time":  DeltaPs,
	"Rpl_semi_sync_master_tx_waits":      DeltaPs,
}

// SlaveStatus not all slave status send to falcon-agent, this is a filter
var SlaveStatus = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

func dataType(k string) string {
	if v, ok := DataType[k]; ok {
		return v
	}
	return Origin
}

// MetaData is the json obj to output
type MetaData struct {
	Metric      string      `json:"metric"`      // metric name
	Endpoint    string      `json:"endpoint"`    // Hostname
	Value       interface{} `json:"value"`       // number or string
	CounterType string      `json:"counterType"` // GAUGE:original value, COUNTER: delta value(ps)
	Tags        string      `json:"tags"`        // port=%d,isSlave=%d,readOnly=%d,type=mysql
	Timestamp   int64       `json:"timestamp"`   // time.Now().Unix()
	Step        int64       `json:"step"`        // Default 60 seconds
}

// NewMetric is the constructor of metric
func NewMetric(conf *common.Config, name string) *MetaData {
	return &MetaData{
		Metric:      name,
		Endpoint:    common.Hostname(conf),
		CounterType: dataType(name),
		Tags:        Tag,
		Timestamp:   time.Now().Unix(),
		Step:        60,
	}
}

// SetValue is just set a value
func (m *MetaData) SetValue(v interface{}) {
	m.Value = v
}

// SetName is just set a value
func (m *MetaData) SetName(name string) {
	m.Metric = name
}

// GetTag can get the tag to output
func GetTag(conf *common.Config) string {
	return fmt.Sprintf(
		"port=%d,isSlave=%d,readOnly=%d,type=mysql",
		conf.DataBase.Port, IsSlave, IsReadOnly)
}

// MySQLAlive checks if mysql can response
func MySQLAlive(conf *common.Config, ok bool) {
	data := NewMetric(conf, "mysql_alive_local")
	data.SetValue(0)
	if ok {
		data.SetValue(1)
	}
	msg, err := SendData(conf, []*MetaData{data})
	if err != nil && err != io.EOF {
		Log.Error("Send alive data failed: %v", err)
		return
	}
	Log.Info("Alive data response %s:%d: %s",
		conf.DataBase.Host, conf.DataBase.Port, string(msg))
}

// GetIsReadOnly get read_only variable of mysql
func GetIsReadOnly(db mysql.Conn) (int, error) {
	row, _, err := db.QueryFirst("select @@read_only")
	if err != nil {
		return -1, err
	}
	return row.Int(0), nil
}
