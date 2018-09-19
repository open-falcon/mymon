# README.md

## Introduction

mymon(MySQL-Monitor) 是Open-Falcon用来监控MySQL数据库运行状态的一个插件，采集包括global status, global variables, slave status以及innodb status等MySQL运行状态信息。

## Installation

```bash
# Build
go get -u github.com/open-falcon/mymon
cd $GOPATH/src/github.com/open-falcon/mymon
make

# Add to crontab
echo '* * * * * cd ${WORKPATH} && ./mymon -c etc/myMon.cfg' > /etc/cron.d/mymon
```

## Configuration

配置文件采用INI标准。 

```ini
[default]
basedir = . # 工作目录
log_dir = ./fixtures # 日志目录，默认日志文件为myMon.log,旧版本有log_file项，如果同时设置了，会优先采用log_file
ignore_file = ./falconignore # 配置忽略的metric项
snapshot_dir = ./snapshot # 保存快照(process, innodb status)的目录
snapshot_day = 10 # 保存快照的时间(日)
log_level  = 5 #  日志级别[RFC5424]
# 0 LevelEmergency
# 1 LevelAlert
# 2 LevelCritical
# 3 LevelError
# 4 LevelWarning
# 5 LevelNotice
# 6 LevelInformational
# 7 LevelDebug
falcon_client=http://127.0.0.1:1988/v1/push # falcon agent连接地址

[mysql]
user=root # 数据库用户名
password=1tIsB1g3rt # 您的数据库密码
host=127.0.0.1 # 数据库连接地址
port=3306 # 数据库端口
```

## Metric

采集的metric信息，请参考./metrics.txt。该文件仅供参考，实际采集信息会根据MySQL版本、配置的不同而变化。

## Contributors

* libin 微信：libin_cc 邮件：libin_dba@xiaomi.com [OLD]
* liuzidong [![Chat on gitter](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/sylzd) 邮件：liuzidong@xiaomi.com [CURRENT]
