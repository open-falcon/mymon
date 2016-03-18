package job

import (
	"github.com/coraldane/mymon/g"
	"github.com/coraldane/mymon/models"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func GlobalStatus(server *g.DBServer, db mysql.Conn) ([]*models.MetaData, error) {
	return mysqlState(server, db, "SHOW /*!50001 GLOBAL */ STATUS")
}

func GlobalVariables(server *g.DBServer, db mysql.Conn) ([]*models.MetaData, error) {
	return mysqlState(server, db, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

func mysqlState(server *g.DBServer, db mysql.Conn, sql string) ([]*models.MetaData, error) {
	rows, _, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	data := make([]*models.MetaData, len(rows))
	i := 0
	for _, row := range rows {
		key_ := row.Str(0)
		v, err := row.Int64Err(1)
		// Ignore non digital value
		if err != nil {
			continue
		}

		data[i] = models.NewMetric(key_, server)
		data[i].SetValue(v)
		i++
	}
	return data[:i], nil
}
