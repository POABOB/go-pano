import redis_service as RS
import threading
import time


HOST = "redis"
PORT = 6379
PWD = "ppaannoo"
# stream对应的key
stream_name = "stream:pano"
# 消费组
consumer_group = "pano_group"

def execute(msgId,data):
    """
    业务处理，不建议多个场景共用一个stream，建议分开，
    如果数据量比较少，通过工厂处理分发
    :param msgId:
    :param data:
    :return:
    """
    print(f'業務執行msgId={msgId} data={data}')
    time.sleep(5)
    return True


def subscribe():
    #基于redis的
    service = RS.RedisService(host=HOST, port=PORT, pwd=PWD)
    print("Subscribe Start")
    r = service.getRedisFromPool()
    stream = RS.RedisStreamService(r, stream_name, consumer_group)
    stream.xgroup_create(id="0")
    while True:
        r = service.getRedisFromPool()
        stream = RS.RedisStreamService(r, stream_name, consumer_group)
        stream.consumer("當前data_id",target=execute, block=0)
        # time.sleep(3)

def test():
    # 初始化
    service = RS.RedisService(host=HOST,port=PORT,pwd=PWD)
    # 池子內獲取連線實例
    r = service.getRedisFromPool()
    
    # 初始化Stream，定義Stram名稱、消費者名稱
    stream =  RS.RedisStreamService(r,stream_name,consumer_group)

    # Init自動創建
    stream.xgroup_create(id="0")
    print("Publish Start")
    for i in range(1, 11):
        stream.xadd({
            "predict_id": i,
            "dir": "dir",
            "filename": "a.jpg",
            "predict_string": "{}",
            "created_at": "2022-10-10",
            "updated_at": "2022-10-10"
        })
        print("test", i)
        time.sleep(1)


if __name__ == '__main__':
    #起一个后台线程执行消费，防止阻塞主线程
    receiveDataThread = threading.Thread(target=subscribe)
    receiveDataThread.start()
    test()



