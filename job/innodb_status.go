package job

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/coraldane/mymon/g"
	"github.com/coraldane/mymon/models"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func InnodbStatus(server *g.DBServer, db mysql.Conn) ([]*models.MetaData, error) {
	status, _, err := db.QueryFirst("SHOW /*!50000 ENGINE */ INNODB STATUS")
	if err != nil {
		return nil, err
	}
	ctn := status.Str(2)
	rows := strings.Split(ctn, "\n")
	return parseInnodbStatus(server, rows)
}

func parseInnodbStatus(server *g.DBServer, rows []string) ([]*models.MetaData, error) {
	var section string
	data := make([]*models.MetaData, 0)
	for _, row := range rows {
		switch {
		case match("^BACKGROUND THREAD$", row):
			section = "BACKGROUND THREAD"
			continue
		case match("^DEAD LOCK ERRORS$", row), match("^LATEST DETECTED DEADLOCK$", row):
			section = "DEAD LOCK ERRORS"
			continue
		case match("^FOREIGN KEY CONSTRAINT ERRORS$", row), match("^LATEST FOREIGN KEY ERROR$", row):
			section = "FOREIGN KEY CONSTRAINT ERRORS"
			continue
		case match("^SEMAPHORES$", row):
			section = "SEMAPHORES"
			continue
		case match("^TRANSACTIONS$", row):
			section = "TRANSACTIONS"
			continue
		case match("^FILE I/O$", row):
			section = "FILE I/O"
			continue
		case match("^INSERT BUFFER AND ADAPTIVE HASH INDEX$", row):
			section = "INSERT BUFFER AND ADAPTIVE HASH INDEX"
			continue
		case match("^LOG$", row):
			section = "LOG"
			continue
		case match("^BUFFER POOL AND MEMORY$", row):
			section = "BUFFER POOL AND MEMORY"
			continue
		case match("^ROW OPERATIONS$", row):
			section = "ROW OPERATIONS"
			continue
		}

		if section == "SEMAPHORES" {
			matches := regexp.MustCompile(`^Mutex spin waits\s+(\d+),\s+rounds\s+(\d+),\s+OS waits\s+(\d+)`).FindStringSubmatch(row)
			if len(matches) == 4 {
				spin_waits, _ := strconv.Atoi(matches[1])
				Innodb_mutex_spin_waits := models.NewMetric("Innodb_mutex_spin_waits", server)
				Innodb_mutex_spin_waits.SetValue(spin_waits)
				data = append(data, Innodb_mutex_spin_waits)

				spin_rounds, _ := strconv.Atoi(matches[2])
				Innodb_mutex_spin_rounds := models.NewMetric("Innodb_mutex_spin_rounds", server)
				Innodb_mutex_spin_rounds.SetValue(spin_rounds)
				data = append(data, Innodb_mutex_spin_rounds)

				os_waits, _ := strconv.Atoi(matches[3])
				Innodb_mutex_os_waits := models.NewMetric("Innodb_mutex_os_waits", server)
				Innodb_mutex_os_waits.SetValue(os_waits)
				data = append(data, Innodb_mutex_os_waits)
			}
		}
	}
	return data, nil
}

func match(pattern, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}
