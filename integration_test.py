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
    message = b'test-message'
    length = await client.enqueue('test-1', message)
    assert length == 1


async def test_length(server):
    client = new_client()
    length = await client.length('test-1')
    assert length == 0
    message = b'test-message'
    length = await client.enqueue('test-1', message)
    assert length == 1
    length = await client.length('test-1')
    assert length == 1


async def test_dequeue(server):
    client = new_client()
    message = b'test-message'
    await client.enqueue('test-1', message)
    value = await client.dequeue('test-1')
    assert value == message


async def test_dequeue_blocks(server):
    client = new_client()
    with pytest.raises(asyncio.TimeoutError):
        await asyncio.wait_for(client.dequeue('test-1'), timeout=0.1)


@pytest.fixture()
def server():
    # We sleep here so that the server has enough time to give up the port between tests.
    time.sleep(0.01)
    server = subprocess.Popen(["build/litemq"], env={"PORT": "42099"})
    # We sleep here to give the server time to start before the tests run.
    time.sleep(0.01)
    yield
    server.kill()


def new_client():
    channel = Channel('127.0.0.1', 42099)
    return LiteMQ(channel)
