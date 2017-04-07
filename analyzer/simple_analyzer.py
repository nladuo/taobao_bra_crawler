#!/usr/bin/env python
# -*- coding:utf-8 -*-

import re
from utils import *
import json
import sys
from datetime import datetime

reload(sys)
sys.setdefaultencoding('utf8')


class SimpleAnalyzer:
    """ 简单的统计工具，根据商品评论的size_info来分析 """

    def __init__(self):
        self.client = init_client()
        self.db = self.client[config['db_name']]
        self.rates = self.db.rates.find({})
        # 初始化数据
        self.sizes = {}
        self.colors = {}
        self.size_details = {}

    def run(self):
        # 使用正则表达式匹配A、B、C、D杯....
        size_pattern = r'[A-K]'
        # 使用正则表达式匹配70A、70B....
        size_detail_pattern = r'[5-9][0-9][A-K]{1}'
        # 只统计以下几种颜色
        color_keys = [u"红色", u"橙色", u"黄色", u"绿色", u"蓝色", u"紫色", u"黑色", u"白色", u"粉色"]

        print "正在统计中， 请耐心等待......"
        before_exec_time = datetime.now()
        for rate in self.rates:
            # print rate['size_info']

            # 统计罩杯尺寸
            sizes = re.findall(size_pattern, rate['size_info'])
            for size in set(sizes):
                if self.sizes.has_key(size):
                    self.sizes[size] += 1
                else:
                    self.sizes[size] = 1

            # 统计罩杯具体尺寸
            size_details = re.findall(size_detail_pattern, rate['size_info'])
            for size_detail in set(size_details):
                if self.size_details.has_key(size_detail):
                    self.size_details[size_detail] += 1
                else:
                    self.size_details[size_detail] = 1

            # 统计罩杯颜色
            for color_key in color_keys:
                if color_key in rate['size_info']:
                    if self.colors.has_key(color_key):
                        self.colors[color_key] += 1
                    else:
                        self.colors[color_key] = 1
        after_exec_time = datetime.now()
        print "共耗时", (after_exec_time - before_exec_time).seconds, "秒"
        self.__close()

    def __close(self):
        """ 关闭数据库 """
        self.client.close()



if __name__ == '__main__':
    analyzer = SimpleAnalyzer()
    analyzer.run()

    bra_data = {
        'basic': analyzer.sizes,
        'color': analyzer.colors,
        'detail': analyzer.size_details
    }

    # 把数据保存到bra.json中
    dat = json.dumps(bra_data, ensure_ascii=False)
    print dat
    f = open('bra.json', 'w')
    f.write(dat)
    f.close()