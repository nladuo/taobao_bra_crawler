# -*- coding:utf-8 -*-
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine, Column, Integer, MetaData, Table, VARCHAR, SmallInteger
import requests

def init_session():
    engine = create_engine('mysql+mysqlconnector://root:root@localhost:3306/taobao')
    create_table(engine)
    DBSession = sessionmaker(bind=engine)
    return DBSession()

def get_body(url):
    return requests.get(url).content

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
        Column('rate_id', Integer, unique=True, nullable=False),
        Column('title', VARCHAR(255), nullable=False),
        Column('is_crawled', VARCHAR(255), nullable=False)
    )


    metadata.create_all(engine)