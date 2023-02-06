#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# @Time    : 2022/4/02 7:35 下午
# @FileName: redisService.py

import redis
from redis.client import *
from redis.sentinel import Sentinel
import json

class RedisService(object):
    """
    直接操作redis
    """
    def __init__(self,host=None,port=None,pwd=None,db=0,max_connections=10):
        """
        redis链接初始化
        :param host: redis实例的地址
        :param port:  端口
        :param pwd: 密码
        :param db:  使用的db
        :param max_connections: 连接池大小
        """
        self.pool = None
        self.host = host
        self.port = port
        self.pwd = pwd
        self.db = db
        self.max_connections = max_connections

    def __getPool(self,flag):
        """
        获取连接池（主从切换以后把flag设置为True，重新初始化连接池）
        :param flag: 为False不用重新获取，为True需要重新获取
        :param max_connections: 连接池大小
        :return:
        """
        try:
            if self.pool is not None:
                return self.pool
            # 缓存连接池   
            self.pool =  redis.ConnectionPool(host=self.host, port=self.port, db=self.db, password=self.pwd, max_connections=self.max_connections)
            return self.pool
        except Exception as e:
            print(e)
            raise e
    def getRedisFromPool(self,flag = False) -> Redis:
        """
        从链接池获取redis链接
        :param flag: 是否重新获取(主从切换后)
        :return:
        """
        return redis.Redis(connection_pool=self.__getPool(flag))
    def getRedis(self) -> Redis:
        """
        直接获取redis链接
        :return:
        """
        if self.pwd:
            return redis.Redis(host=self.host, port=self.port, db=self.db, password=self.pwd)
        return redis.Redis(host=self.host, port=self.port, db=self.db)

class RedisStreamService(object):
    """
    redis stream封装
    """
    def __init__(self,redis,stream_name,consumer_group):
        """
        :param redis:redis的客户端
        :param stream_name:stream的名称
        :param consumer_group:消费组名称
        """
        self.redis = redis
        self.stream_name = stream_name
        self.consumer_group = consumer_group

    def __streamInit(self,data,id="0",maxlen=20000,target = None):
        """
        初始化stream，上线之前手动调用下即可，不用在项目里调用
        :param data: 业务数据  测试数据
        :param id: 0 从开始消费, $ 从创建以后新进来的开始消费
        :param maxlen: 队里最大长度
        :return:
        """
        if not self.redis.exists(self.stream_name):
            """
            不存在消费者，直接创建消费者和消费组
            """
            # self.redis.xadd(self.stream_name, self.__dataWrap(data), maxlen=maxlen)
            self.xgroup_create(id)
            # self.redis.xadd(self.stream_name, self.__dataWrap(data), maxlen=maxlen)
        # self.consumer("robot",count=2,target=target)
    def xgroup_create(self,id="0",maxlen=20000):
        """
        创建消费组
        :param id:
        :return:
        """
        try:
            rsp = self.redis.xinfo_groups(self.stream_name)
            for item in rsp:
                if self.consumer_group == item['name'].decode('utf8'):
                    print(f'xgroup:{self.consumer_group}已存在, 无需再次创建')
                    return
        except Exception as msg:
            pass

        group = self.redis.xgroup_create(self.stream_name, self.consumer_group, id=id, mkstream=True)
        if not group:
            raise Exception(f'创建消费者组:{self.consumer_group}异常')

    def __dataWrap(self,data) -> dict:
        """
        包装数据
        :param data:
        :return:
        """
        return  {"data":json.dumps(data)}
    def xack(self,msgId):
        """
        ack
        :param msgId:
        :return:
        """
        self.redis.xack( self.stream_name,  self.consumer_group, msgId)
    def xdel(self,msgId):
        """
        del
        :param msgId:
        :return:
        """
        self.redis.xdel( self.stream_name, msgId)

    def __getData(self,data):
        """
        从消息流中获取业务数据
        :param item:
        :return:
        """
        if not data or not data[0]:
            return None, None
        msgId = str(data[0], 'utf-8')
        data = {str(key, 'utf-8'): str(val, 'utf-8') for key, val in data[1].items()}
        return msgId, data["data"]
    def xadd(self,data):
        """
        新增数据
        :param stream_name:
        :param data:
        :return:
        """
        self.redis.xadd( self.stream_name, self.__dataWrap(data))
    # 業務處理
    def consumer(self,consumer_name,id=">",block=60, count=5,target=None):
        """
        消费数据
        :param consumer_name: 消费者名称，建议传递ip
        :param id: 从哪开始消费
        :param block: 无消息阻塞时间，毫秒，默认60秒，在60秒内有消息直接消费
        :param count: 消费多少条，默认1
        :param target: 业务处理回调方法
        :return:
        """
        # block 0 时阻塞等待, 其他数值表示读取超时时间
        streams = {self.stream_name: id}
        rst = self.redis.xreadgroup( self.consumer_group, consumer_name, streams, block=block, count=count)
        print(f'消费到的数据 {rst}')
        if not rst or not rst[0] or not rst[0][1]:
           return None
        # 遍历获取到的列表信息（可以消费多条，根据count）
        for item in rst[0][1]:
           try:
               #解析数据
               msgId, data = self.__getData(item)
               """
               执行回调函数target，成功后ack
               """
               if target and target(msgId,data):
                   # 将处理完成的消息标记，类似于kafka的offset
                   self.redis.xack( self.stream_name,  self.consumer_group, msgId)
                   self.redis.xdel( self.stream_name, msgId)
           except Exception as e:
               # 消费失败，下次从头消费(消费成功的都已经提交ack了，可以先不处理，以后再处理)
               print("consumer is error:",e)
