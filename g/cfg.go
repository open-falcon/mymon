package g

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type DBServer struct {
	Endpoint string `json:"endpoint"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Passwd   string `json:"passwd"`
}

func (this *DBServer) String() string {
	return fmt.Sprintf("<Endpoint: %s, Host: %s, Port: %d, User: %s, Passwd: %s>",
		this.Endpoint, this.Host, this.Port, this.User, this.Passwd)
}

type GlobalConfig struct {
	LogLevel       string      `json:"log_level"`
	Interval       int         `json:"interval"`
	ConnectTimeout int         `json:"connect_timeout"`
	FalconClient   string      `json:"falcon_client"`
	DBServerList   []*DBServer `json:"db_server_list"`
}

var (
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("config file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read config file %s fail %s", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse config file %s fail %s", cfg, err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
	return nil
}
