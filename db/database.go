package db

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/coraldane/mymon/g"
	"github.com/toolkits/logger"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ormLock = new(sync.RWMutex)
)

func InitDatabase() {
	// set default database
	if g.Config().LogLevel == "debug" {
		orm.Debug = true
	}

	maxIdle := g.Config().MaxIdle
	for index, server := range g.Config().DBServerList {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?loc=Local&parseTime=true",
			server.User, server.Passwd, server.Host, server.Port)
		fmt.Println(dsn)
		if 0 == index {
			orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxIdle)
		}
		orm.RegisterDataBase(server.Endpoint, "mysql", dsn, maxIdle, maxIdle)
	}
}

func NewOrmWithAlias(alias string) orm.Ormer {
	ormLock.RLock()
	defer ormLock.RUnlock()

	o := orm.NewOrm()
	o.Using(alias)
	return o
}

func QueryFirst(alias, strSql string, args ...interface{}) (orm.Params, error) {
	var maps []orm.Params
	num, err := NewOrmWithAlias(alias).Raw(strSql, args...).Values(&maps)
	if nil != err {
		logger.Errorln(num, err)
		return nil, err
	}
	if num > 0 {
		return maps[0], err
	}
	return nil, err
}

func QueryRows(alias, strSql string, args ...interface{}) ([]orm.Params, error) {
	var maps []orm.Params
	num, err := NewOrmWithAlias(alias).Raw(strSql, args...).Values(&maps)
	if nil != err {
		logger.Errorln(num, err)
		return nil, err
	}
	return maps, err
}
