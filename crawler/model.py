# -*- coding:utf-8 -*-
from sqlalchemy import Column, Integer, MetaData, Table, VARCHAR, BOOLEAN
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
    size_info = Column(VARCHAR(100))
    rate_content = Column(VARCHAR(100))

    def __init__(self, size_info, rate_content):
        self.size_info = size_info
        self.rate_content = rate_content

