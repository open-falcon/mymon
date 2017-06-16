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
	"fmt"
	"os"
	"time"
)

const (
	TIME_OUT = 30

	ORIGIN   = "GAUGE"
	DELTA_PS = "COUNTER"
	DELTA    = ""
)

// COUNTER: Speed per second
// GAUGE: Original, DEFAULT
var DataType = map[string]string{
	"Innodb_buffer_pool_reads":         DELTA_PS,
	"Innodb_buffer_pool_read_requests": DELTA_PS,
	"Innodb_compress_time":             DELTA_PS,
	"Innodb_data_fsyncs":               DELTA_PS,
	"Innodb_data_read":                 DELTA_PS,
	"Innodb_data_reads":                DELTA_PS,
	"Innodb_data_writes":               DELTA_PS,
	"Innodb_data_written":              DELTA_PS,
	"Innodb_last_checkpoint_at":        DELTA_PS,
	"Innodb_log_flushed_up_to":         DELTA_PS,
	"Innodb_log_sequence_number":       DELTA_PS,
	"Innodb_mutex_os_waits":            DELTA_PS,
	"Innodb_mutex_spin_rounds":         DELTA_PS,
	"Innodb_mutex_spin_waits":          DELTA_PS,
	"Innodb_pages_flushed_up_to":       DELTA_PS,
	"Innodb_rows_deleted":              DELTA_PS,
	"Innodb_rows_inserted":             DELTA_PS,
	"Innodb_rows_locked":               DELTA_PS,
	"Innodb_rows_modified":             DELTA_PS,
	"Innodb_rows_read":                 DELTA_PS,
	"Innodb_rows_updated":              DELTA_PS,
	"Innodb_row_lock_time":             DELTA_PS,
	"Innodb_row_lock_waits":            DELTA_PS,
	"Innodb_uncompress_time":           DELTA_PS,

	"Binlog_event_count": DELTA_PS,
	"Binlog_number":      DELTA_PS,
	"Slave_count":        DELTA_PS,

	"Com_admin_commands":        DELTA_PS,
	"Com_assign_to_keycache":    DELTA_PS,
	"Com_alter_db":              DELTA_PS,
	"Com_alter_db_upgrade":      DELTA_PS,
	"Com_alter_event":           DELTA_PS,
	"Com_alter_function":        DELTA_PS,
	"Com_alter_procedure":       DELTA_PS,
	"Com_alter_server":          DELTA_PS,
	"Com_alter_table":           DELTA_PS,
	"Com_alter_tablespace":      DELTA_PS,
	"Com_analyze":               DELTA_PS,
	"Com_begin":                 DELTA_PS,
	"Com_binlog":                DELTA_PS,
	"Com_call_procedure":        DELTA_PS,
	"Com_change_db":             DELTA_PS,
	"Com_change_master":         DELTA_PS,
	"Com_check":                 DELTA_PS,
	"Com_checksum":              DELTA_PS,
	"Com_commit":                DELTA_PS,
	"Com_create_db":             DELTA_PS,
	"Com_create_event":          DELTA_PS,
	"Com_create_function":       DELTA_PS,
	"Com_create_index":          DELTA_PS,
	"Com_create_procedure":      DELTA_PS,
	"Com_create_server":         DELTA_PS,
	"Com_create_table":          DELTA_PS,
	"Com_create_trigger":        DELTA_PS,
	"Com_create_udf":            DELTA_PS,
	"Com_create_user":           DELTA_PS,
	"Com_create_view":           DELTA_PS,
	"Com_dealloc_sql":           DELTA_PS,
	"Com_delete":                DELTA_PS,
	"Com_delete_multi":          DELTA_PS,
	"Com_do":                    DELTA_PS,
	"Com_drop_db":               DELTA_PS,
	"Com_drop_event":            DELTA_PS,
	"Com_drop_function":         DELTA_PS,
	"Com_drop_index":            DELTA_PS,
	"Com_drop_procedure":        DELTA_PS,
	"Com_drop_server":           DELTA_PS,
	"Com_drop_table":            DELTA_PS,
	"Com_drop_trigger":          DELTA_PS,
	"Com_drop_user":             DELTA_PS,
	"Com_drop_view":             DELTA_PS,
	"Com_empty_query":           DELTA_PS,
	"Com_execute_sql":           DELTA_PS,
	"Com_flush":                 DELTA_PS,
	"Com_grant":                 DELTA_PS,
	"Com_ha_close":              DELTA_PS,
	"Com_ha_open":               DELTA_PS,
	"Com_ha_read":               DELTA_PS,
	"Com_help":                  DELTA_PS,
	"Com_insert":                DELTA_PS,
	"Com_insert_select":         DELTA_PS,
	"Com_install_plugin":        DELTA_PS,
	"Com_kill":                  DELTA_PS,
	"Com_load":                  DELTA_PS,
	"Com_lock_tables":           DELTA_PS,
	"Com_optimize":              DELTA_PS,
	"Com_preload_keys":          DELTA_PS,
	"Com_prepare_sql":           DELTA_PS,
	"Com_purge":                 DELTA_PS,
	"Com_purge_before_date":     DELTA_PS,
	"Com_release_savepoint":     DELTA_PS,
	"Com_rename_table":          DELTA_PS,
	"Com_rename_user":           DELTA_PS,
	"Com_repair":                DELTA_PS,
	"Com_replace":               DELTA_PS,
	"Com_replace_select":        DELTA_PS,
	"Com_reset":                 DELTA_PS,
	"Com_resignal":              DELTA_PS,
	"Com_revoke":                DELTA_PS,
	"Com_revoke_all":            DELTA_PS,
	"Com_rollback":              DELTA_PS,
	"Com_rollback_to_savepoint": DELTA_PS,
	"Com_savepoint":             DELTA_PS,
	"Com_select":                DELTA_PS,
	"Com_set_option":            DELTA_PS,
	"Com_signal":                DELTA_PS,
	"Com_show_authors":          DELTA_PS,
	"Com_show_binlog_events":    DELTA_PS,
	"Com_show_binlogs":          DELTA_PS,
	"Com_show_charsets":         DELTA_PS,
	"Com_show_collations":       DELTA_PS,
	"Com_show_contributors":     DELTA_PS,
	"Com_show_create_db":        DELTA_PS,
	"Com_show_create_event":     DELTA_PS,
	"Com_show_create_func":      DELTA_PS,
	"Com_show_create_proc":      DELTA_PS,
	"Com_show_create_table":     DELTA_PS,
	"Com_show_create_trigger":   DELTA_PS,
	"Com_show_databases":        DELTA_PS,
	"Com_show_engine_logs":      DELTA_PS,
	"Com_show_engine_mutex":     DELTA_PS,
	"Com_show_engine_status":    DELTA_PS,
	"Com_show_events":           DELTA_PS,
	"Com_show_errors":           DELTA_PS,
	"Com_show_fields":           DELTA_PS,
	"Com_show_function_status":  DELTA_PS,
	"Com_show_grants":           DELTA_PS,
	"Com_show_keys":             DELTA_PS,
	"Com_show_master_status":    DELTA_PS,
	"Com_show_open_tables":      DELTA_PS,
	"Com_show_plugins":          DELTA_PS,
	"Com_show_privileges":       DELTA_PS,
	"Com_show_procedure_status": DELTA_PS,
	"Com_show_processlist":      DELTA_PS,
	"Com_show_profile":          DELTA_PS,
	"Com_show_profiles":         DELTA_PS,
	"Com_show_relaylog_events":  DELTA_PS,
	"Com_show_slave_hosts":      DELTA_PS,
	"Com_show_slave_status":     DELTA_PS,
	"Com_show_status":           DELTA_PS,
	"Com_show_storage_engines":  DELTA_PS,
	"Com_show_table_status":     DELTA_PS,
	"Com_show_tables":           DELTA_PS,
	"Com_show_triggers":         DELTA_PS,
	"Com_show_variables":        DELTA_PS,
	"Com_show_warnings":         DELTA_PS,
	"Com_slave_start":           DELTA_PS,
	"Com_slave_stop":            DELTA_PS,
	"Com_stmt_close":            DELTA_PS,
	"Com_stmt_execute":          DELTA_PS,
	"Com_stmt_fetch":            DELTA_PS,
	"Com_stmt_prepare":          DELTA_PS,
	"Com_stmt_reprepare":        DELTA_PS,
	"Com_stmt_reset":            DELTA_PS,
	"Com_stmt_send_long_data":   DELTA_PS,
	"Com_truncate":              DELTA_PS,
	"Com_uninstall_plugin":      DELTA_PS,
	"Com_unlock_tables":         DELTA_PS,
	"Com_update":                DELTA_PS,
	"Com_update_multi":          DELTA_PS,
	"Com_xa_commit":             DELTA_PS,
	"Com_xa_end":                DELTA_PS,
	"Com_xa_prepare":            DELTA_PS,
	"Com_xa_recover":            DELTA_PS,
	"Com_xa_rollback":           DELTA_PS,
	"Com_xa_start":              DELTA_PS,

	"Aborted_clients":            DELTA_PS,
	"Aborted_connects":           DELTA_PS,
	"Access_denied_errors":       DELTA_PS,
	"Binlog_bytes_written":       DELTA_PS,
	"Binlog_cache_disk_use":      DELTA_PS,
	"Binlog_cache_use":           DELTA_PS,
	"Binlog_stmt_cache_disk_use": DELTA_PS,
	"Binlog_stmt_cache_use":      DELTA_PS,
	"Bytes_received":             DELTA_PS,
	"Bytes_sent":                 DELTA_PS,
	"Connections":                DELTA_PS,
	"Created_tmp_disk_tables":    DELTA_PS,
	"Created_tmp_files":          DELTA_PS,
	"Created_tmp_tables":         DELTA_PS,
	"Handler_delete":             DELTA_PS,
	"Handler_read_first":         DELTA_PS,
	"Handler_read_key":           DELTA_PS,
	"Handler_read_last":          DELTA_PS,
	"Handler_read_next":          DELTA_PS,
	"Handler_read_prev":          DELTA_PS,
	"Handler_read_rnd":           DELTA_PS,
	"Handler_read_rnd_next":      DELTA_PS,
	"Handler_update":             DELTA_PS,
	"Handler_write":              DELTA_PS,
	"Opened_files":               DELTA_PS,
	"Opened_tables":              DELTA_PS,
	"Opened_table_definitions":   DELTA_PS,
	"Qcache_hits":                DELTA_PS,
	"Qcache_inserts":             DELTA_PS,
	"Qcache_lowmem_prunes":       DELTA_PS,
	"Qcache_not_cached":          DELTA_PS,
	"Queries":                    DELTA_PS,
	"Questions":                  DELTA_PS,
	"Select_full_join":           DELTA_PS,
	"Select_full_range_join":     DELTA_PS,
	"Select_range_check":         DELTA_PS,
	"Select_scan":                DELTA_PS,
	"Slow_queries":               DELTA_PS,
	"Sort_merge_passes":          DELTA_PS,
	"Sort_range":                 DELTA_PS,
	"Sort_rows":                  DELTA_PS,
	"Sort_scan":                  DELTA_PS,
	"Table_locks_immediate":      DELTA_PS,
	"Table_locks_waited":         DELTA_PS,
	"Threads_created":            DELTA_PS,
}

type MysqlIns struct {
	Host string
	Port int
	Tag  string
}

func dataType(key_ string) string {
	if v, ok := DataType[key_]; ok {
		return v
	}
	return ORIGIN
}

type MetaData struct {
	Metric      string      `json:"metric"`      //key
	Endpoint    string      `json:"endpoint"`    //hostname
	Value       interface{} `json:"value"`       // number or string
	CounterType string      `json:"counterType"` // GAUGE  原值   COUNTER 差值(ps)
	Tags        string      `json:"tags"`        // port=3306,k=v
	Timestamp   int64       `json:"timestamp"`
	Step        int64       `json:"step"`
}

func (m *MetaData) String() string {
	s := fmt.Sprintf("MetaData Metric:%s Endpoint:%s Value:%v CounterType:%s Tags:%s Timestamp:%d Step:%d",
		m.Metric, m.Endpoint, m.Value, m.CounterType, m.Tags, m.Timestamp, m.Step)
	return s
}

func NewMetric(name string) *MetaData {
	return &MetaData{
		Metric:      name,
		Endpoint:    hostname(),
		CounterType: dataType(name),
		Tags:        fmt.Sprintf("port=%d", cfg.Port),
		Timestamp:   time.Now().Unix(),
		Step:        60,
	}
}

func hostname() string {
	host := cfg.Endpoint
	if host != "" {
		return host
	}
	host, err := os.Hostname()
	if err != nil {
		host = cfg.Host
	}
	return host
}

func (m *MetaData) SetValue(v interface{}) {
	m.Value = v
}
