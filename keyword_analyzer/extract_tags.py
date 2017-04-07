#!/usr/bin/env python
# -*- coding:utf-8 -*-
""" 提取关键词 """
from datetime import datetime
import jieba
import jieba.analyse
import sys
import cPickle as pickle

reload(sys)
sys.setdefaultencoding('utf8')


if __name__ == '__main__':

    print "Reading the corpus..."
    content = open("./assets/corpus.txt", 'r').read()
    print "Finished reading."

    print "Extracting tags..."
    t1 = datetime.now()
    tags = jieba.analyse.extract_tags(content, topK=100, withWeight=True)
    print "Finished extraction in", (datetime.now() - t1).seconds, "second(s)."

    for tag in tags:
        print tag[0], "\t", tag[1]

    # 保存tags
    with open("./assets/tags.pickle", "w") as f:
        pickle.dump(tags, f)

