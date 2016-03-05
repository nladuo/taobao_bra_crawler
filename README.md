# taobao_bra_crawler

## 说明
淘宝文胸商品评论内容爬取与简单分析，支持单机爬取和分布式爬取。<br>
效果展示：[http://vps.kalen25115.cn/bra](http://vps.kalen25115.cn/bra/)

## 注意
由于使用github.com/qiniu/iconv库的原因，目前好像不支持go1.6版本。测试使用的是go1.5.3版本。

## 项目状态
测试中，目前都不是很稳定。

## Dependency
``` go
go get github.com/nladuo/go-webcrawler  # 自己写的一个简单的分布式的爬虫框架，正在慢慢完善
go get github.com/bitly/go-simplejson   # 复杂json操作的库
go get github.com/Go-SQL-Driver/MySQL   # mysql驱动
go get github.com/qiniu/iconv           # 编码转换库,需要安装gcc
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
## 爬虫部署流程
### 1. 爬取商品的记录
crawler文件下的每个go文件都有下面数据库的配置信息，根据自己的情况修改一下mysql的配置，并创建数据库taobao。
``` go
const (
        DB_USER   = "root"      //用户名
        DB_PASSWD = "root"      //密码
        DB_HOST   = "localhost" //地址
        DB_PORT   = "3306"      //端口号
        DBNAME    = "taobao"    //数据库名称
)
```
修改过后，直接go run就可以了，有关商品记录的代码是在crawler/bra_item_crawler.go或者crawler/bra_item_crawler_without_proxy.go中。一个是使用代理的，一个是不用代理的，第一个因为数据量不是太大，所以不用代理也可以。<br><br>
注：爬虫程序不会自动停止，可以看着情况把它停掉，根据商品条目的数量来定，一般两千条以上就差不多，商品的条目放在了taobao下面的bras表中。

### 2. 爬取文胸商品的评论
爬取完商品用不了太长时间，如果不用代理的话，大概几十秒就可以爬到两千多。用代理的话，也不需要太久。<br><br>
接下来就要开始爬取商品下的评论了，每个商品都有可能有一千条评论，这么算下来两千个商品记录最多应该有最多200万的评论，这么多数据，需要上代理了。（如果不用代理的话，大概只能爬到一万多评论）。
#### 单机部署
一台机器爬的话，还是和之前一样，就修改一下mysql的配置，直接用go run就可以了。
代码在crawler/bra_rates_crawler/single目录中。

#### 集群部署（目前不稳定,容易卡住不动）
分布式部署需要首先配置一下zookeeper，<a href="http://zookeeper.apache.org/doc/r3.4.6/zookeeperStarted.html">点击这里查看</a>。<br><br>
把zookeeper启动起来之后，来看一下crawler/bra_rates_crawler/single目录下的config.json
``` json
{
    "AppName": "taobao-bra-crawler",
    "IsMaster": true,
    "ThreadNum": 200,
    "LockerTimeout": 10,
    "ZkTimeOut": 600,
    "ZkHosts": [
        "127.0.0.1:2181"
    ]
}
```
这里主要修改一下ZkHosts，这里是个json数组，如果是zk集群，按照相应的填写就可以了。<br><br>
还有就是IsMaster，只有一台机器可以设置为true。其他的机器修改为false，就可以了。<br><br>
运行方式:
``` go
go run bra_rates_crawler.go config.json
```
用go run再加上config.json参数就okay。

##代理ip说明
代码中使用的代理ip来自免费的代理ip服务器，（http://www.89ip.cn 这个网站）。免费代理ip不太稳定，而且可能不是高匿名的，可以换收费的试试。

## 爬取后的分析
当爬到一定的商品数据之后，就可以进行一定的分析了，简单的分析测试在simple_analyser.go中，修改一下mysql的配置，直接go run就好了。
``` go
go run simple_analyser.go >> bra.json
```
把运行的结果存入bra.json中，然后替换掉web_display文件夹的bra.json文件。把web_display文件夹放到一个web服务器的根目录中，比如说apache，输入localhost/web_display 就可以看到效果了。

## 可能存在的bug
#### 跑着跑着程序不动了，很久stdout都没有输出log信息
此时可以停止掉爬虫，如果是单机模式运行的话把addBaseTasks()这行注释掉，然后再运行，会继续工作。如果是集群模式，如果是master结点，把配置文件中的IsMaster改成false就可以，如果不是master，直接重启就okay。
#### 出现大量too many open files的信息
在linux下，因为socket也是文件，并发量比较高，需要重新设置一下最大文件打开的数量。
``` shell
sudo su             #切换root
ulimit -n 50000     #设置一个程序可以打开5万个文件
```


## License
MIT