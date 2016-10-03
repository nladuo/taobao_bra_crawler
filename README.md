# taobao_bra_crawler
a taobao web crawler just for fun.

## 说明
淘宝文胸商品评论内容爬取与简单分析。<br>
测试环境：腾讯云主机一台，系统是ubuntu-14.04<br>
数据库： mongodb
效果展示：[http://nladuo.github.io/bra](http://nladuo.github.io/bra)

## 数据下载


## 项目状态
正在使用python重写中。。。

## 爬虫部署
### 修改配置文件
``` python
# -*- coding:utf-8 -*-
config = {
    'timeout' : 3,
    'db_user': '',       # 无密码
    'db_pass': '',
    'db_host': 'localhost',
    'db_port': 27017,
    'db_name': 'taobao',
    'use_tor_proxy': False,
    'tor_proxy_port': 9050
}
```
如果有被禁IP的情况可以使用tor代理，将config['use_tor_proxy']设置为True，具体方法见[python中使用tor代理](http://nladuo.github.io/2016/07/17/python%E4%B8%AD%E4%BD%BF%E7%94%A8tor%E4%BB%A3%E7%90%86/)
### 运行爬虫
``` shell
python crawler/item_crawler.py      #爬文胸的商品信息
python crawler/rate_crawler.py      #爬文胸的评论信息
python crawler/simple_analyzer.py   #统计数据
```
### 运行网页显示
``` shell
cd data_visualization
npm install
npm run dev
node build/build.js # 生成dist
```
