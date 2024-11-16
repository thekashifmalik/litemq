import asyncio
import grpc
from grpclib.client import Channel

from service_pb2 import EnqueueRequest
from service_pb2 import QueueID
from service_pb2 import Nothing
from service_grpc import QueueServiceStub


async def main():
    async with Channel('127.0.0.1', 42090) as channel:
        greeter = QueueServiceStub(channel)

        reply = await greeter.Health(Nothing())
        print(reply)

        reply = await greeter.Enqueue(EnqueueRequest(queue='hello', data=b'asdasd'))
        print(reply)

        reply = await greeter.Dequeue(QueueID(queue='hello'))
        print(reply.data)
        reply = await greeter.Dequeue(QueueID(queue='hello'))
        print(reply.data)


if __name__ == '__main__':
    asyncio.run(main())
