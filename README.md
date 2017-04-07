# taobao_bra_crawler
a taobao web crawler just for fun.

## 说明
淘宝文胸商品评论内容爬取与简单分析。

## 部署环境
测试环境：腾讯云主机一台<br>
操作系统：ubuntu-14.04<br>
数据库： mongodb<br>
效果展示：[http://nladuo.github.io/bra](http://nladuo.github.io/bra)

## 商品评论数据
### 下载地址
链接: [https://pan.baidu.com/s/1bpbuZLX](https://pan.baidu.com/s/1bpbuZLX) 密码: kvyp

### 导入数据
``` bash
mongoimport -d taobao -c rates  --file ./rates.dat
```

## 爬虫部署
### 安装依赖
``` bash
pip install -r crawler/requirements.txt
```
### 修改配置文件
``` python
config = {
    'timeout' : 3,
    'db_user': '',
    'db_pass': '',
    'db_host': 'localhost',
    'db_port': 27017,
    'db_name': 'taobao',
    'use_tor_proxy': False,
    'tor_proxy_port': 9050
}
```
一般爬取速度不快不会有禁IP的情况。如果有被禁IP的情况可以使用tor代理，将config['use_tor_proxy']设置为True，具体方法见[python中使用tor代理](http://nladuo.github.io/2016/07/17/python%E4%B8%AD%E4%BD%BF%E7%94%A8tor%E4%BB%A3%E7%90%86/)
### 运行爬虫
``` bash
python crawler/item_crawler.py      # 爬文胸的商品信息
python crawler/rate_crawler.py      # 爬文胸的评论信息
```
### 数据分析
#### 简单统计与可视化展示
1. 统计数据
``` sh
cd simple_analyzer
python simple_analyzer.py               # 简单统计
cp bra.json data_visualization/static/  # 拷贝统计结果
```
2. 运行网页显示
``` sh
cd data_visualization
npm install     # 安装依赖
npm run dev     # 进行调试
npm run build   # 生成dist
```
#### 关键词分析
编写中...