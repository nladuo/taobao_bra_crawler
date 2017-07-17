# -*- coding:utf-8 -*-
import pymongo
import requests
import requesocks
from model import FailedUrl
from config import *


# 初始化mongodb客户端
def init_client():
    client = pymongo.MongoClient(config['db_host'], config['db_port'])
    if len(config['db_user']) != 0:
        admin = client[config['db_name']]
        admin.authenticate(config['db_user'], config['db_pass'])
    return client


# 根据是否使用tor代理，来获取http客户端
def get_http_client():
    if config['use_tor_proxy']:
        session = requesocks.session()
        session.proxies = {'http': 'socks5://127.0.0.1:%d' % config['tor_proxy_port'],
                           'https': 'socks5://127.0.0.1:%d' % config['tor_proxy_port']}
        return session
    else:
        return requests.session()


# 发送get请求
def get_body(url):
    retry_times = 0
    client = get_http_client()
    while retry_times < 3:
        try:
            content = client.get(url, timeout=config['timeout']).content
            return content
        except:
            retry_times += 1
    return ''


# 把失败的url添加到数据库
def add_failed_url(db, url):
    collection = db.failed_urls
    if collection.find({'url': url}).count() == 0:
        collection.insert(FailedUrl(url).dict())

