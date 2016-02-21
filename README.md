# taobao_bra_crawler

## 说明
淘宝文胸商品评论内容爬取与简单分析，支持单机爬取和集群爬取。<br>
效果展示：[http://vps.kalen25115.cn/bra/](http://vps.kalen25115.cn/bra/)

## 状态
测试中。。。

## Dependency
``` go
go get github.com/nladuo/go-webcrawler  # 自己写的一个简单的分布式的爬虫框架，正在慢慢完善
go get github.com/bitly/go-simplejson   # 复杂json操作的库
go get github.com/Go-SQL-Driver/MySQL   # mysql驱动
go get github.com/qiniu/iconv           # 编码转换库
```
## Installation
### Mac or Linux
``` shell
git clone https://github.com/nladuo/taobao_bra_crawler.git
cd taobao_bra_crawler
chmod +x dependency.sh
./dependency.sh
```

### Windows
``` shell
git clone https://github.com/nladuo/taobao_bra_crawler.git
cd taobao_bra_crawler
dependency.bat
```
## 代码部分改动
```
db, err := gorm.Open("mysql", "root:root@/taobao?charset=utf8&parseTime=True&loc=Local")
```
代码用到了mysql数据库，当然你也可以使用sqlite3或者postgreSQL,只要更改驱动就okay。
## 单机模式
```
cd single
go run bra_crawler.go
```
## 集群模式
### 搭建zk环境
Check out the zookeeper configuration <a href="http://zookeeper.apache.org/doc/r3.4.6/zookeeperStarted.html">here</a>.
### Master结点
更改配置文件，主要修改ZkHosts，根据自己的ip以及设置的通信端口来定。
``` json
{
    "AppName": "taobao-bra-crawler",
    "IsMaster": true,
    "ThreadNum": 2000,
    "LockerTimeout": 10,
    "ZkTimeOut": 600,
    "ZkHosts": [
        "192.168.1.102:2181" //zookeeper配置，这里是zk单机运行
    ]
}
```
### Worker结点
更改配置文件，把IsMaster改为false
``` json
{
    "AppName": "taobao-bra-crawler",
    "IsMaster": false,
    "ThreadNum": 2000,
    "LockerTimeout": 10,
    "ZkTimeOut": 600,
    "ZkHosts": [
        "192.168.1.102:2181" //zookeeper配置，这里是zk单机运行
    ]
}
```
### 运行
```
go run bra_crawler.go config.json
```
## License
MIT