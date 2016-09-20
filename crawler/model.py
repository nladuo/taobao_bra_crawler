# -*- coding:utf-8 -*-
from sqlalchemy import Column, Integer, VARCHAR, BOOLEAN
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class Item(Base):
    __tablename__ = 'items'

    id = Column(Integer, primary_key=True, autoincrement=True)
    item_id = Column(VARCHAR(50), unique=True)
    seller_id = Column(VARCHAR(50))
    title = Column(VARCHAR(100))
    is_crawled = Column(BOOLEAN)

    def __init__(self, item_id, seller_id, title, is_crawled):
        self.item_id = item_id
        self.seller_id = seller_id
        self.title = title
        self.is_crawled = is_crawled


class Rate(Base):
    __tablename__ = 'rates'

    id = Column(Integer, primary_key=True, autoincrement=True)
    rate_id = Column(VARCHAR(100), unique=True)
    size_info = Column(VARCHAR(100))
    rate_content = Column(VARCHAR(100))

    def __init__(self, rate_id, size_info, rate_content):
        self.rate_id = rate_id
        self.size_info = size_info
        self.rate_content = rate_content


class FailedUrl(Base):
    __tablename__ = 'failed_urls'

    id = Column(Integer, primary_key=True, autoincrement=True)
    url = Column(VARCHAR(255), unique=True)

    def __init__(self, url):
        self.url = url

