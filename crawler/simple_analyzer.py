# -*- coding:utf-8 -*-

import re
from utils import *
import json
import sys

reload(sys)
sys.setdefaultencoding('utf8')

class SimpleAnalyzer:

    def __init__(self):
        self.client = init_client()
        self.db = self.client[config['db_name']]
        self.rates = self.db.rates.find({})
        self.sizes = {}
        self.colors = {}
        self.size_details = {}



    def run(self):
        # 使用正则表达式匹配A、B、C、D杯....
        size_pattern = r'[A-K]'
        # 使用正则表达式匹配70A、70B....
        size_detail_pattern = r'[5-9][0-9][A-K]{1}'
        color_keys = [u"红色", u"橙色", u"黄色", u"绿色", u"蓝色", u"紫色", u"黑色", u"白色", u"粉色"]
        for rate in self.rates:
            print rate['size_info']
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




analyzer = SimpleAnalyzer()
analyzer.run()
bra_data = {
    'basic': analyzer.sizes,
    'color': analyzer.colors,
    'detail': analyzer.size_details
}
dat = json.dumps(bra_data, ensure_ascii=False)
print dat
f = open('bra.json', 'w')
f.write(dat)
f.close()