# -*- coding:utf-8 -*-
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine, Column, Integer, MetaData, Table, VARCHAR, SmallInteger
import requests
from model import FailedUrl

def init_session():
    engine = create_engine('mysql+mysqlconnector://root:root@localhost:3306/taobao')
    create_table(engine)
    DBSession = sessionmaker(bind=engine)
    return DBSession()

def get_body(url):
    retry_times = 0
    while retry_times < 3:
        try:
            content = requests.get(url, timeout=3).content
            return content
        except:
            retry_times += 1
    return ''


def add_failed_url(session, url):
    if session.query(FailedUrl).filter(FailedUrl.url == url).count() == 0:
        session.add(FailedUrl(url))
        session.commit()

def create_table(engine):
    metadata = MetaData()
    Table(
        "items", metadata,
        Column('id', Integer, primary_key=True, autoincrement=True),
        Column('item_id', VARCHAR(100), unique=True, nullable=False),
        Column('seller_id', VARCHAR(100), nullable=False),
        Column('title', VARCHAR(255), nullable=False),
        Column('is_crawled', SmallInteger, nullable=False)
    )

    Table(
        "rates", metadata,
        Column('id', Integer, primary_key=True, autoincrement=True),
        Column('rate_id', VARCHAR(100), unique=True, nullable=False),
        Column('size_info', VARCHAR(255), nullable=False),
        Column('rate_content', VARCHAR(255), nullable=False)
    )

    Table(
        "failed_urls", metadata,
        Column('id', Integer, primary_key=True, autoincrement=True),
        Column('url', VARCHAR(255), unique=True, nullable=False),
    )

    metadata.create_all(engine)