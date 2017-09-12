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

// MySQL Performance Monitor(For open-falcon)
// Write by Li Bin<libin_dba@xiaomi.com>
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	goconf "github.com/akrennmair/goconf"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Cfg struct {
	LogFile      string
	LogLevel     int
	FalconClient string
	Endpoint     string

	User string
	Pass string
	Host string
	Port int
}

var cfg Cfg

func init() {
	var cfgFile string
	flag.StringVar(&cfgFile, "c", "myMon.cfg", "myMon configure file")
	flag.Parse()

	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			log.WithField("cfg", cfgFile).Fatalf("myMon config file does not exists: %v", err)
		}
	}

	if err := cfg.readConf(cfgFile); err != nil {
		log.Fatalf("Read configure file failed: %v", err)
	}

	// Init log file
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(cfg.LogLevel))

	if cfg.LogFile != "" {
		f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			log.SetOutput(f)
			return
		}
	}
	log.SetOutput(os.Stderr)
}

func (conf *Cfg) readConf(file string) error {
	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		return err
	}

	conf.LogFile, err = c.GetString("default", "log_file")
	if err != nil {
		return err
	}

	conf.LogLevel, err = c.GetInt("default", "log_level")
	if err != nil {
		return err
	}

	conf.FalconClient, err = c.GetString("default", "falcon_client")
	if err != nil {
		return err
	}

	conf.Endpoint, err = c.GetString("default", "endpoint")
	if err != nil {
		return err
	}

	conf.User, err = c.GetString("mysql", "user")
	if err != nil {
		return err
	}

	conf.Pass, err = c.GetString("mysql", "password")
	if err != nil {
		return err
	}

	conf.Host, err = c.GetString("mysql", "host")
	if err != nil {
		return err
	}

	conf.Port, err = c.GetInt("mysql", "port")
	return err
}

func timeout() {
	time.AfterFunc(TIME_OUT*time.Second, func() {
		MysqlAlive(nil, false)
		log.Error("Execute timeout")
		os.Exit(1)
	})
}

func MysqlAlive(m *MysqlIns, ok bool) {
	data := NewMetric("mysql_alive_local")
	if ok {
		data.SetValue(1)
	} else {
		data.SetValue(0)
	}
	msg, err := sendData([]*MetaData{data})
	if err != nil {
		log.Errorf("Send alive data failed: %v", err)
		return
	}
	if m != nil {
		log.Infof("Alive data response %s: %s", m.String(), string(msg))
	}
}

func FetchData(m *MysqlIns) (err error) {
	defer func() {
		MysqlAlive(m, err == nil)
	}()

	db := mysql.New("tcp", "", fmt.Sprintf("%s:%d", m.Host, m.Port),
		cfg.User, cfg.Pass)
	db.SetTimeout(500 * time.Millisecond)
	if err = db.Connect(); err != nil {
		return
	}
	defer db.Close()

	data := make([]*MetaData, 0)
	globalStatus, err := GlobalStatus(m, db)
	if err != nil {
		return
	}
	data = append(data, globalStatus...)

	globalVars, err := GlobalVariables(m, db)
	if err != nil {
		return
	}
	data = append(data, globalVars...)

	innodbState, err := innodbStatus(m, db)
	if err != nil {
		return
	}
	data = append(data, innodbState...)

	slaveState, err := slaveStatus(m, db)
	if err != nil {
		return
	}
	data = append(data, slaveState...)

	msg, err := sendData(data)
	if err != nil {
		return
	}
	log.Infof("Send response %s: %s", m.String(), string(msg))
	return
}

func (m *MysqlIns) String() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func main() {
	log.Info("MySQL Monitor for falcon")
	go timeout()

	err := FetchData(&MysqlIns{
		Host: cfg.Host,
		Port: cfg.Port,
		Tag:  fmt.Sprintf("port=%d", cfg.Port),
	})
	if err != nil {
		log.Error(err)
	}
}
