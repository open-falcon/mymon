## Introduction

mymon(MySQL-Monitor) -- MySQL数据库运行状态数据采集脚本，采集包括global status, global variables, slave status等。

由github.com/open-falcon/mymon 修改而来，
支持配置多个MySQL实例
## Installation

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/coraldane
cd $GOPATH/src/github.com/coraldane
git clone https://github.com/coraldane/mymon.git

cd mymon
go get ./...
control build


```

## Configuration

```
{
    "log_level": "debug",
    "interval": 60,
    "falcon_client": "http://127.0.0.1:1988/v1/push",
    "db_server_list": [
    {
        "endpoint": "机器名1",
        "host": "host地址1",
        "port": 3306,
        "user": "$user",
        "passwd": "$pass"
    },
    {
        "endpoint": "机器名2",
        "host": "host地址2",
        "port": 3306,
        "user": "$user",
        "passwd": "$pass"
    }
    ]
}
```

## MySQL metrics

请参考./metrics.txt，其中的内容，仅供参考，根据MySQL的版本、配置不同，采集到的metrics也有差别。


## Contributors

 - libin  微信：libin_cc  邮件：libin_dba@xiaomi.com
 - coraldane 邮件：coraldane@163.com
