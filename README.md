## Introduction

mymon(MySQL-Monitor) -- MySQL数据库运行状态数据采集脚本，采集包括global status, global variables, slave status等。

## Installation

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/mymon.git

cd mymon
go get ./...
go build -o mymon

echo '* * * * * cd $GOPATH/src/github.com/open-falcon/mymon && ./mymon -c etc/mon.cfg' > /etc/cron.d/mymon

```

## Configuration

```
    [default]
    log_file=mymon.log # 日志路径和文件名
    # Panic 0
    # Fatal 1
    # Error 2
    # Warn 3
    # Info 4
    # Debug 5
    log_level=4 # 日志级别

    falcon_client=http://127.0.0.1:1988/v1/push # falcon agent连接地址

    #自定义endpoint
    endpoint=127.0.0.1 #若不设置则使用OS的hostname

    [mysql]
    user=root # 数据库用户名
    password= # 数据库密码
    host=127.0.0.1 # 数据库连接地址
    port=3306 # 数据库端口
```

## MySQL metrics

请参考./metrics.txt，其中的内容，仅供参考，根据MySQL的版本、配置不同，采集到的metrics也有差别。


## Contributors

 - libin  微信：libin_cc  邮件：libin_dba@xiaomi.com

