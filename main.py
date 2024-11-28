import asyncio
import grpc
from grpclib.client import Channel

from litemq.client import LiteMQ


async def main():
    async with Channel('127.0.0.1', 42090) as channel:
        client = LiteMQ(channel)
        print('sending health')
        await client.health()
        print('sending enq')
        length = await client.enqueue('hello', b'asdasd')
        print(length)
        print('sending deq')
        data = await client.dequeue('hello')
        print(data)


if __name__ == '__main__':
    asyncio.run(main())
    # with grpc.insecure_channel('127.0.0.1:42090') as channel:
    #     stub = QueueServiceStub(channel)
    #     stub.Enqueue(EnqueueRequest(queue='lol', data=b'lolol'))
