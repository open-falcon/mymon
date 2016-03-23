package job

import (
	"fmt"
	"github.com/coraldane/mymon/db"
	"github.com/coraldane/mymon/g"
	"github.com/coraldane/mymon/models"
	"strconv"
)

var SlaveStatusToSend = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_Log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

func SlaveStatus(server *g.DBServer) ([]*models.MetaData, error) {
	isSlave := models.NewMetric("Is_slave", server)

	row, err := db.QueryFirst(g.Hostname(server), "SHOW SLAVE STATUS")
	if err != nil {
		return nil, err
	}

	// be master
	if row == nil || 0 == len(row) {
		isSlave.SetValue(0)
		return []*models.MetaData{isSlave}, nil
	}

	// be slave
	isSlave.SetValue(1)

	data := make([]*models.MetaData, len(SlaveStatusToSend))
	for i, s := range SlaveStatusToSend {
		data[i] = models.NewMetric(s, server)
		switch s {
		case "Slave_SQL_Running", "Slave_IO_Running":
			data[i].SetValue(0)
			v := fmt.Sprintf("%v", row[s])
			if v == "Yes" {
				data[i].SetValue(1)
			}
		default:
			v, err := strconv.Atoi(fmt.Sprintf("%v", row[s]))
			if err != nil {
				data[i].SetValue(-1)
			} else {
				data[i].SetValue(v)
			}
		}
	}
	return append(data, isSlave), nil
}
