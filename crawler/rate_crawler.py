#!/usr/bin/env python
# -*- coding:utf-8 -*-

from model import *
from utils import *
import json

class RateCrawler:
    def __init__(self):
        self.session = init_session()
        self.items = self.session.query(Item)

    def run(self):
        for item in self.items:
            base_url = "https://rate.tmall.com/list_detail_rate.htm?itemId=%s&sellerId=%s&currentPage=%d&pageSize=1000000"
            url = base_url % (item.item_id, item.seller_id, 1)
            body = "{" + get_body(url).decode("gbk") + "}"



        self.__close()

    def __parse_page_num(self, body):
        items = []
        try:
            data = json.loads(body)
        except:
            return []
        item_list = data['listItem']
        if len(item_list) == 0:
            return []
        for _item in item_list:
            item = Item(_item['item_id'], _item['userId'], _item['title'], False)
            items.append(item)
        return items

    def __parse_rates(self, body):
        items = []
        try:
            data = json.loads(body)
        except:
            return []
        item_list = data['listItem']
        if len(item_list) == 0:
            return []
        for _item in item_list:
            item = Item(_item['item_id'], _item['userId'], _item['title'], False)
            items.append(item)
        return items

    def __add(self, rate):
        if self.session.query(Rate).filter(Rate.rate_id == rate.rate_id).count() == 0:
            self.session.add(rate)
            self.session.commit()

    def __update_item(self, item):
        pass

    def __close(self):
        self.session.close()

url = "https://rate.tmall.com/list_detail_rate.htm?itemId=536964054394&sellerId=2262562480&currentPage=1&pageSize=1000000"

body = "{" + get_body(url).decode("gbk") + "}"

data = json.loads(body)

fileHandle = open ( 'test.txt', 'w' )
fileHandle.write(json.dumps(data['rateDetail']['rateList'][0]))
