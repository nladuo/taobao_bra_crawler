#!/usr/bin/env python
# -*- coding:utf-8 -*-
""" 生成词云图片 """
from wordcloud import WordCloud, ImageColorGenerator
import cPickle as pickle
import numpy as np
from PIL import Image

if __name__ == "__main__":

    # 读取词频
    with open("./assets/tags.pickle", "r") as f:
        tags = pickle.load(f)
        frequencies = {}
        for tag in tags:
            frequencies[tag[0]] = int(10000 * tag[1])

    wc = WordCloud(font_path='./assets/simhei.ttf',  # 设置字体
                   background_color="black",  # 背景颜色
                   max_words=100,  # 词云显示的最大词数
                   max_font_size=80,  # 字体最大值
                   random_state=42)

    wc.generate_from_frequencies(frequencies)

    # 颜色转换
    rainbow_coloring = np.array(Image.open("./assets/rainbow.jpg"))
    image_colors = ImageColorGenerator(rainbow_coloring)
    wc.recolor(color_func=image_colors)

    # 保存图片
    wc.to_file("./assets/word_cloud.png")
    print "saved at ./assets/word_cloud.png"
