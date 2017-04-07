#!/usr/bin/env python
# -*- coding:utf-8 -*-

from gevent import monkey; monkey.patch_all()
import gevent
from gevent import queue
import json
import time
import sys
sys.path.append("../")
from lib.model import *
from lib.utils import *

reload(sys)
sys.setdefaultencoding('utf8')


class RateCrawler:
    """ 根据文胸商品, 爬取评论 """

    def __init__(self):
        self.client = init_client()
        self.db = self.client[config['db_name']]
        self.collection = self.db.rates
        self.collection.ensure_index('rate_id', unique=True)
        self.items = self.db.items.find({'is_crawled': False})

    def run(self):
        items = []
        # 先把数据读到内存
        for item in self.items:
            items.append(Item(
                item['item_id'],
                item['seller_id'],
                item['title'],
                False
            ))
            pass

        for item in items:
            base_url = "https://rate.tmall.com/list_detail_rate.htm?itemId=%s&sellerId=%s&currentPage=%d&pageSize=1000000"
            url = base_url % (item.item_id, item.seller_id, 1)
            try:
                # 这里返回的数据不是纯json，需要在两边加上{}
                body = "{" + get_body(url).decode("gbk") + "}"
                if len(body) == 2:
                    add_failed_url(self.db, url)
                    continue
            except:
                add_failed_url(self.db, url)
                continue

            # 获取评论页数
            page_num = self.__parse_page_num(body)
            print item.title, ' ', item.item_id, '--------->' , page_num, \
                time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(time.time()))

            # 使用gevent并发爬取，把数据存在queue里
            tasks = []
            q = gevent.queue.Queue()
            for i in range(1, page_num+1):
                url = base_url % (item.item_id, item.seller_id, i)
                tasks.append(gevent.spawn(self.__async_get_rates, url, q))
            gevent.joinall(tasks)
            print "adding data of item:%s" % item.item_id,  time.strftime("%Y-%m-%d %H:%M:%S",
                                                                          time.localtime(time.time()))
            # 逐个添加到数据库
            while not q.empty():
                body = q.get()
                if len(body) == 2:
                    add_failed_url(self.db, url)
                    continue
                rates = self.__parse_rates(body)
                self.__add_rates(rates)
            # 把item的is_crawled设为1
            self.__update_item(item)
            # time.sleep(30) # 睡眠30秒
        self.__close()

    def __async_get_rates(self, url, q):
        """ 异步发送get请求 """
        try:
            body = "{" + get_body(url).decode("gbk") + "}"
            q.put(body)
        except:
            add_failed_url(self.db, url)
        print url

    def __parse_page_num(self, body):
        """ 解析商品的评论页数 """
        try:
            data = json.loads(body)
            page_num = data['rateDetail']['paginator']['lastPage']
            return page_num
        except:
            return 0

    def __parse_rates(self, body):
        """ 解析商品的评论 """
        rates = []
        try:
            data = json.loads(body)
        except:
            return []
        rate_list = data['rateDetail']['rateList']
        if len(rate_list) == 0:
            return []
        for _rate in rate_list:
            rate = Rate(_rate['id'], _rate['auctionSku'], _rate['rateContent'])
            rates.append(rate)
        return rates

    def __add_rates(self, rates):
        """ 添加商品评论 """
        for rate in rates:
            try:
                self.collection.insert(rate.dict())
            except:
                pass

    def __update_item(self, item):
        """ 把当前商品设置为：已经爬取过 """
        self.db.items.update({'item_id': item.item_id}, {
            '$set': {'is_crawled': True},
        })

    def __close(self):
        """ 关闭数据库 """
        self.client.close()


if __name__ == '__main__':
    crawler = RateCrawler()
    crawler.run()

