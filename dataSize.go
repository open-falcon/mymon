package main

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"strconv"
)

func getDataSize(m *MysqlIns, db mysql.Conn) (*MetaData, error) {
	sql := "SELECT (SUM(DATA_LENGTH)+SUM(INDEX_LENGTH)) AS data_size FROM INFORMATION_SCHEMA.TABLES"

	row, _, err := db.QueryFirst(sql)
	if err != nil {
		return nil, err
	}
	sizeStr := row.Str(0)
	size, _ := strconv.Atoi(sizeStr)

	dataSize := NewMetric("Data_size")
	dataSize.SetValue(size)

	return dataSize, nil
}
