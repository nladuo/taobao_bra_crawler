#!/usr/bin/env python
# -*- coding:utf-8 -*-

import json
import sys
sys.path.append("../")
from lib.model import *
from lib.utils import *
from lib.config import *

reload(sys)
sys.setdefaultencoding('utf8')


class ItemCrawler:
    """ 爬取淘宝文胸商品记录 """

    def __init__(self):
        self.client = init_client()
        self.db = self.client[config['db_name']]
        self.collection = self.db.items

    def run(self):
        urls = []
        for i in range(1, 100 + 1):
            urls.append("http://s.m.taobao.com/search?q=文胸&m=api4h5&page=" + str(i))
            urls.append("http://s.m.taobao.com/search?q=胸罩&m=api4h5&page=" + str(i))
            urls.append("http://s.m.taobao.com/search?q=bra&m=api4h5&page=" + str(i))

        for url in urls:
            print url
            body = get_body(url)
            if len(body) == 0:
                continue
            items = self.__parse(body)
            self.__add_items(items)

        self.__close()

    def __parse(self, body):
        """ 解析商品记录 """
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
        """ 添加商品记录到数据库 """
        for item in items:
            if self.collection.find({'item_id': item.item_id}).count() == 0:
                self.collection.insert(item.dict())

    def __close(self):
        """ 关闭数据库 """
        self.client.close()


if __name__ == '__main__':
    crawler = ItemCrawler()
    crawler.run()

