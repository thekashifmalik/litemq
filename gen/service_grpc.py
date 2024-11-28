import abc
import typing
import grpclib.const
import grpclib.client
if typing.TYPE_CHECKING:
    import grpclib.server
from . import service_pb2

class QueueServiceBase(abc.ABC):

    @abc.abstractmethod
    async def Enqueue(self, stream: 'grpclib.server.Stream[service_pb2.EnqueueRequest, service_pb2.QueueLength]') -> None:
        pass

    @abc.abstractmethod
    async def Dequeue(self, stream: 'grpclib.server.Stream[service_pb2.QueueID, service_pb2.DequeueResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Purge(self, stream: 'grpclib.server.Stream[service_pb2.QueueID, service_pb2.QueueLength]') -> None:
        pass

    @abc.abstractmethod
    async def Length(self, stream: 'grpclib.server.Stream[service_pb2.QueueID, service_pb2.QueueLength]') -> None:
        pass

    @abc.abstractmethod
    async def Health(self, stream: 'grpclib.server.Stream[service_pb2.Nothing, service_pb2.Nothing]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {'/QueueService/Enqueue': grpclib.const.Handler(self.Enqueue, grpclib.const.Cardinality.UNARY_UNARY, service_pb2.EnqueueRequest, service_pb2.QueueLength), '/QueueService/Dequeue': grpclib.const.Handler(self.Dequeue, grpclib.const.Cardinality.UNARY_UNARY, service_pb2.QueueID, service_pb2.DequeueResponse), '/QueueService/Purge': grpclib.const.Handler(self.Purge, grpclib.const.Cardinality.UNARY_UNARY, service_pb2.QueueID, service_pb2.QueueLength), '/QueueService/Length': grpclib.const.Handler(self.Length, grpclib.const.Cardinality.UNARY_UNARY, service_pb2.QueueID, service_pb2.QueueLength), '/QueueService/Health': grpclib.const.Handler(self.Health, grpclib.const.Cardinality.UNARY_UNARY, service_pb2.Nothing, service_pb2.Nothing)}

class QueueServiceStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.Enqueue = grpclib.client.UnaryUnaryMethod(channel, '/QueueService/Enqueue', service_pb2.EnqueueRequest, service_pb2.QueueLength)
        self.Dequeue = grpclib.client.UnaryUnaryMethod(channel, '/QueueService/Dequeue', service_pb2.QueueID, service_pb2.DequeueResponse)
        self.Purge = grpclib.client.UnaryUnaryMethod(channel, '/QueueService/Purge', service_pb2.QueueID, service_pb2.QueueLength)
        self.Length = grpclib.client.UnaryUnaryMethod(channel, '/QueueService/Length', service_pb2.QueueID, service_pb2.QueueLength)
        self.Health = grpclib.client.UnaryUnaryMethod(channel, '/QueueService/Health', service_pb2.Nothing, service_pb2.Nothing)