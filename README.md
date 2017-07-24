# taobao_bra_crawler
a taobao web crawler just for fun.

## 说明
淘宝文胸商品评论内容爬取与简单分析。

## 商品评论数据
### 下载地址
链接: [https://pan.baidu.com/s/1bpbuZLX](https://pan.baidu.com/s/1bpbuZLX) 密码: kvyp

### 导入数据
``` bash
mongoimport -d taobao -c rates  --file ./rates.dat
```

## 爬虫部署
互联网时代的网站富于变化，爬虫今天可能正常明天可能就不能用了，如果爬虫无法使用请通过百度云盘链接导入数据。
### 部署环境
测试环境：腾讯云主机一台<br>
操作系统：ubuntu-14.04<br>
数据库： mongodb<br>

### 安装依赖
``` bash
pip install -r requirements.txt
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
说明：一般的爬取速度不会有禁IP的情况。
### 运行爬虫
``` bash
python crawler/item_crawler.py      # 爬文胸的商品信息
python crawler/rate_crawler.py      # 爬文胸的评论信息
```

## 数据处理
### 简单统计与可视化展示
#### 1. 运行脚本
``` sh
cd simple_analyzer
python simple_analyzer.py               # 简单统计
cp bra.json data_visualization/static/  # 拷贝统计结果
```
#### 2. 运行网页显示
``` sh
cd data_visualization
npm install     # 安装依赖
npm run dev     # 进行调试
npm run build   # 生成dist
```
#### 效果展示
见: [http://nladuo.github.io/bra](http://nladuo.github.io/bra)

### 关键词分析
#### 运行脚本
``` sh
cd keyword_analyzer
python create_corpus.py     # 1.加载评论信息
python extract_tags.py      # 2.提取关键词(20分钟左右, 可以直接用我的模型进行第三步)
python create_wordcloud.py  # 3.生成词云图片
```
#### 效果
![word_cloud](./keyword_analyzer/assets/word_cloud1.png)

#### 参考
- [Python pytagcloud 中文分词 生成标签云 系列（一）](http://www.cnblogs.com/Yiutto/p/5998262.html)
- [利用pandas+python制作100G亚马逊用户评论数据词云](http://www.jianshu.com/p/c862130f322d)

## LICENSE
MIT