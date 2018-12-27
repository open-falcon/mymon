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

### ignore项
ignore项配置，是用来屏蔽之前在falcon中设好的报警项，会将原有的metric更改名称之后上传，使原有的报警策略不再生效。由于falcon中的屏蔽策略，只能屏蔽endpoint级别，所以在mymon中的ignore功能是帮助提高了报警屏蔽粒度，而非忽略该metric的上报。

### 同步延迟

关于同步延迟检测的metric有两个: `Seconds_Behind_Master`、`Heartbeats_Behind_Master`。

`Seconds_Behind_Master`是MySQL`SHOW SLAVE STATUS`输出的状态变量。由于低版本的MySQL还不支持HEARTBEAT_EVENT，在低版本的MySQL中该状态可能会由于IO线程假死导致测量不准确，因此mymon增加了`Heartbeats_Behind_Master`。它依赖于`pt-heartbeat`，统计基于`pt-heartbeat`生成的mysql.heartbeat表中的ts字段值与从库当前时间差。如果未配置`pt-heartbeat`，则该项上报-1值。

关于pt-heartbeat的配置使用，链接如下：
https://www.percona.com/doc/percona-toolkit/LATEST/pt-heartbeat.html


## Contributors

* libin 微信：libin_cc 邮件：libin_dba@xiaomi.com [OLD]
* liuzidong [![Chat on gitter](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/sylzd) 邮件：liuzidong@xiaomi.com [CURRENT]