package cron

import (
	"bytes"
	"encoding/json"
	"github.com/coraldane/mymon/g"
	"github.com/coraldane/mymon/job"
	"github.com/coraldane/mymon/models"
	"github.com/toolkits/logger"
	"io/ioutil"
	"net/http"
)

func FetchData(server *g.DBServer) (err error) {
	defer func() {
		MysqlAlive(server, err == nil)
	}()

	data := make([]*models.MetaData, 0)
	globalStatus, err := job.GlobalStatus(server)
	if err != nil {
		logger.Errorln("get GlobalStatus error", err)
		return
	}
	data = append(data, globalStatus...)

	globalVars, err := job.GlobalVariables(server)
	if err != nil {
		logger.Errorln("get GlobalVariables error", err)
		return
	}
	data = append(data, globalVars...)

	innodbState, err := job.InnodbStatus(server)
	if err != nil {
		logger.Errorln("get InnodbStatus error", err)
		return
	}
	data = append(data, innodbState...)

	slaveState, err := job.SlaveStatus(server)
	if err != nil {
		logger.Errorln("get SlaveStatus error", err)
		return
	}
	data = append(data, slaveState...)

	msg, err := sendData(data)
	if err != nil {
		logger.Errorln("sendData error", err)
		return
	}
	logger.Info("Send response %s: %s", server.String(), string(msg))
	return
}

func MysqlAlive(server *g.DBServer, ok bool) {
	data := models.NewMetric("mysql_alive_local", server)
	if ok {
		data.SetValue(1)
	}
	msg, err := sendData([]*models.MetaData{data})
	if err != nil {
		logger.Error("Send alive data failed: %v", err)
		return
	}
	logger.Info("Alive data response %s: %s", server.String(), string(msg))
}

func sendData(data []*models.MetaData) ([]byte, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	strUrl := g.Config().FalconClient
	logger.Debug("Send to %s, size: %d", strUrl, len(data))
	for _, m := range data {
		logger.Debug("%s", m)
	}

	res, err := http.Post(strUrl, "Content-Type: application/json", bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
