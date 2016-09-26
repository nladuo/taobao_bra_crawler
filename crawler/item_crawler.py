#!/usr/bin/env python
# -*- coding:utf-8 -*-

from model import *
from utils import *
from config import *
import json
import sys

reload(sys)
sys.setdefaultencoding('utf8')

class ItemCrawler:

    def __init__(self):
        self.client = init_client()
        self.db = self.client[config['db_name']]
        self.collection = self.db.items

    def run(self):
        urls = []
        for i in range(1, 100 + 1):
            urls.append("http://s.m.taobao.com/search?q=文胸&m=api4h5&page=" + str(i))
            urls.append("http://s.m.taobao.com/search?q=胸罩&m=api4h5&page=" + str(i))

        for url in urls:
            print url
            body = get_body(url)
            if len(body) == 0:
                add_failed_url(self.db, url)
                continue
            items = self.__parse(body)
            self.__add_items(items)

        self.__close()


    def __parse(self, body):
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



    def __add_items(self, items):
        for item in items:
            if self.collection.find({'item_id': item.item_id}).count() == 0:
                self.collection.insert(item.dict())

    def __close(self):
        self.client.close()



crawler = ItemCrawler()
crawler.run()


