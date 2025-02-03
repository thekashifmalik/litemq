from grpclib.client import Channel

from .gen.service_pb2 import EnqueueRequest
from .gen.service_pb2 import QueueID
from .gen.service_pb2 import Nothing
from .gen.service_grpc import LiteMQStub


class LiteMQ:
    """
    Example usage:

        from grpclib.client import Channel
        from litemq.client import LiteMQ

        channel = Channel('localhost', 42069)
        client = LiteMQ(channel)
    """

    def __init__(self, channel: Channel):
        self.channel = channel
        self.stub = LiteMQStub(channel)

    async def enqueue(self, queue: str, data: bytes):
        request = EnqueueRequest(queue=queue, data=data)
        response = await self.stub.Enqueue(request)
        return response.count

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

    async def flush(self):
        return await self.stub.Flush(Nothing())
