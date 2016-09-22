#!/usr/bin/env python
# -*- coding:utf-8 -*-

from gevent import monkey; monkey.patch_all()
import gevent
from gevent import queue
from model import *
from utils import *
import json
import time
import sys

reload(sys)
sys.setdefaultencoding('utf8')

class RateCrawler:
    def __init__(self):
        self.session = init_session()
        self.items = self.session.query(Item).filter(Item.is_crawled == False).all()

    def run(self):
        for item in self.items:
            base_url = "https://rate.tmall.com/list_detail_rate.htm?itemId=%s&sellerId=%s&currentPage=%d&pageSize=1000000"
            url = base_url % (item.item_id, item.seller_id, 1)
            try:
                body = "{" + get_body(url).decode("gbk") + "}"
                if len(body) == 2:
                    add_failed_url(self.session, url)
                    continue
            except:
                add_failed_url(self.session, url)
                continue

            page_num = self.__parse_page_num(body)
            print item.title, ' ', item.item_id, '--------->' , page_num, time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(time.time()))
            tasks = []
            q = gevent.queue.Queue()
            for i in range(1, page_num+1):
                url = base_url % (item.item_id, item.seller_id, i)
                tasks.append(gevent.spawn(self.__async_get_rates, url, q))
            gevent.joinall(tasks)
            print "adding data of item:%s" % item.item_id,  time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(time.time()))
            while not q.empty():
                body = q.get()
                if len(body) == 2:
                    add_failed_url(self.session, url)
                    continue
                rates = self.__parse_rates(body)
                self.__add_rates(rates)
            self.__update_item(item) # 把item的is_crawled设为1
        self.__close()

    def __async_get_rates(self, url, q):
        try:
            body = "{" + get_body(url).decode("gbk") + "}"
            q.put(body)
        except:
            add_failed_url(self.session, url)
        print url


    def __parse_page_num(self, body):
        try:
            data = json.loads(body)
            page_num = data['rateDetail']['paginator']['lastPage']
            return page_num
        except:
            return 0


    def __parse_rates(self, body):
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
        for rate in rates:
            if self.session.query(Rate).filter(Rate.rate_id == rate.rate_id).count() == 0:
                self.session.add(rate)
        self.session.commit()

    def __update_item(self, item):
        self.session.query(Item).filter(Item.id == item.id).update({
            'is_crawled': True
        })
        self.session.commit()

    def __close(self):
        self.session.close()

crawler = RateCrawler()
crawler.run()

