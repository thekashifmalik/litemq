import asyncio
import grpc
from grpclib.client import Channel

from gen.service_pb2 import EnqueueRequest
from gen.service_pb2 import QueueID
from gen.service_pb2 import Nothing
from gen.service_grpc import QueueServiceStub


class HTTPQ:
    def __init__(self, channel: Channel):
        self.stub = QueueServiceStub(channel)

    async def enqueue(self, queue: str, data: bytes):
        request = EnqueueRequest(queue=queue, data=data)
        return await self.stub.Enqueue(request)

    async def dequeue(self, queue: str):
        request = QueueID(queue=queue)
        response = await self.stub.Dequeue(request)
        return response.data

    async def length(self, queue: str):
        request = QueueID(queue=queue)
        response = await self.stub.Length(request)
        return response.count

    async def purge(self, queue: str):
        request = QueueID(queue=queue)
        response = await self.stub.Purge(request)
        return response.count

    async def health(self):
        return await self.stub.Health(Nothing())



async def main():
    async with Channel('127.0.0.1', 42090) as channel:
        client = HTTPQ(channel)
        print('sending health')
        await client.health()
        print('sending enq')
        await client.enqueue('hello', b'asdasd')
        print('sending deq')
        data = await client.dequeue('hello')
        print(data)


if __name__ == '__main__':
    asyncio.run(main())
    # with grpc.insecure_channel('127.0.0.1:42090') as channel:
    #     stub = QueueServiceStub(channel)
    #     stub.Enqueue(EnqueueRequest(queue='lol', data=b'lolol'))
