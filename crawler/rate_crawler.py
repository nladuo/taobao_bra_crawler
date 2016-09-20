#!/usr/bin/env python
# -*- coding:utf-8 -*-

from model import *
from utils import *
import json

url = "https://rate.tmall.com/list_detail_rate.htm?itemId=536964054394&sellerId=2262562480&currentPage=1&pageSize=1000000"

body = "{" + get_body(url).decode("gbk") + "}"

data = json.loads(body)

print data['rateDetail']['rateList'][0]