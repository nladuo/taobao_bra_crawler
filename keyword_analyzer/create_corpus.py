#!/usr/bin/env python
# -*- coding:utf-8 -*-
""" 提取评论数据到本地txt """
import sys
sys.path.append("../")
from lib.utils import *
from lib.config import *


reload(sys)
sys.setdefaultencoding('utf8')


def get_rates():
    client = init_client()
    db = client[config['db_name']]
    rates = db.rates.find({})
    return rates

if __name__ == '__main__':
    rates = get_rates()
    print "Adding to corpus.txt..."
    f = open('./assets/corpus.txt', 'a')
    for rate in rates:
        f.write(rate['rate_content'])
        f.write(" ")
    f.close()
    print "Done."
