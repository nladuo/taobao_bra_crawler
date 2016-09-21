# -*- coding:utf-8 -*-
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine, Column, Integer, MetaData, Table, VARCHAR, SmallInteger
import requests
import requesocks
from model import FailedUrl
from config import *

def init_session():
    engine = create_engine('%s+%s://%s:%s@%s:%s/%s'
                           % (config['db_type'], config['db_driver'], config['db_user'],
                              config['db_pass'], config['db_host'],
                              config['db_port'], config['db_name'] ) )
    create_table(engine)
    DBSession = sessionmaker(bind=engine)
    return DBSession()

def get_http_client():
    if config['use_tor_proxy']:
        session = requesocks.session()
        session.proxies = {'http': 'socks5://127.0.0.1:%d' % config['tor_proxy_port'],
                           'https': 'socks5://127.0.0.1:%d' % config['tor_proxy_port']}
        return session
    else:
        return requests.session()

def get_body(url):
    retry_times = 0
    client = get_http_client()
    while retry_times < 3:
        try:
            content = client.get(url, timeout=config['timeout']).content
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