import time
import subprocess
import asyncio
import pytest
import grpc
from grpclib.client import Channel

from litemq.client import LiteMQ


async def test_health(server):
    client = new_client()
    await client.health()


async def test_enqueue(server):
    client = new_client()
    assert await client.enqueue('test', b'message-1') == 1
    assert await client.enqueue('test', b'message-2') == 2


async def test_length(server):
    client = new_client()
    assert await client.length('test') == 0
    await client.enqueue('test', b'message')
    assert await client.length('test') == 1


async def test_dequeue(server):
    client = new_client()
    message = b'message'
    await client.enqueue('test', message)
    assert await client.dequeue('test') == message


async def test_dequeue_blocks(server):
    client = new_client()
    with pytest.raises(asyncio.TimeoutError):
        await asyncio.wait_for(client.dequeue('test-1'), timeout=0.1)


async def test_dequeue_order(server):
    client = new_client()
    message_1 = b'message-1'
    message_2 = b'message-2'
    await client.enqueue('test', message_1)
    await client.enqueue('test', message_2)
    assert await client.dequeue('test') == message_1
    assert await client.dequeue('test') == message_2


@pytest.fixture()
def server():
    # Sleep here so the server has enough time to give up the port between tests.
    time.sleep(0.01)
    server = subprocess.Popen(["build/litemq"], env={"PORT": "42099"})
    # Sleep here to give the server time to start before the tests run.
    time.sleep(0.01)
    yield
    server.kill()


def new_client():
    channel = Channel('127.0.0.1', 42099)
    return LiteMQ(channel)