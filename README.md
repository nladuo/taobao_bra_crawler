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

## License
MIT