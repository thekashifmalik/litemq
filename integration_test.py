import time
import subprocess
import asyncio
import pytest
from grpclib.client import Channel

from litemq.client import LiteMQ


async def test_server(server):
    assert True


async def test_health(server):
    client = new_client(server)
    await client.health()


async def test_enqueue_single(server):
    client = new_client(server)
    assert await client.enqueue('test', b'message-1') == 1


async def test_enqueue_multiple(server):
    client = new_client(server)
    assert await client.enqueue('test', b'message-1') == 1
    assert await client.enqueue('test', b'message-2') == 2


async def test_enqueue_different_queues(server):
    client = new_client(server)
    assert await client.enqueue('test-1', b'message-1') == 1
    assert await client.enqueue('test-2', b'message-1') == 1


async def test_enqueue_different_queues_multiple(server):
    client = new_client(server)
    assert await client.enqueue('test-1', b'message-1') == 1
    assert await client.enqueue('test-1', b'message-2') == 2
    assert await client.enqueue('test-2', b'message-1') == 1
    assert await client.enqueue('test-2', b'message-2') == 2


async def test_length(server):
    client = new_client(server)
    assert await client.length('test') == 0
    await client.enqueue('test', b'message')
    assert await client.length('test') == 1


async def test_dequeue_single(server):
    client = new_client(server)
    await client.enqueue('test', b'message')
    assert await client.dequeue('test') == b'message'


async def test_dequeue_multiple(server):
    client = new_client(server)
    await client.enqueue('test', b'message-1')
    await client.enqueue('test', b'message-2')
    assert await client.dequeue('test') == b'message-1'
    assert await client.dequeue('test') == b'message-2'


async def test_dequeue_different_queues(server):
    client = new_client(server)
    await client.enqueue('test-1', b'message-1')
    await client.enqueue('test-2', b'message-2')
    assert await client.dequeue('test-1') == b'message-1'
    assert await client.dequeue('test-2') == b'message-2'


async def test_dequeue_different_queues_multiple(server):
    client = new_client(server)
    await client.enqueue('test-1', b'message-1')
    await client.enqueue('test-1', b'message-2')
    await client.enqueue('test-2', b'message-1')
    await client.enqueue('test-2', b'message-2')
    assert await client.dequeue('test-1') == b'message-1'
    assert await client.dequeue('test-1') == b'message-2'
    assert await client.dequeue('test-2') == b'message-1'
    assert await client.dequeue('test-2') == b'message-2'


async def test_dequeue_blocks(server):
    client = new_client(server)
    with pytest.raises(asyncio.TimeoutError):
        await asyncio.wait_for(client.dequeue('test-1'), timeout=0.1)


async def test_purge(server):
    client = new_client(server)
    await client.enqueue('test', b'message')
    assert await client.length('test') == 1
    assert await client.purge('test') == 1
    assert await client.length('test') == 0


@pytest.fixture(params=['go', 'rust'])
def server(request):
    # Sleep here so the server has enough time to give up the port between tests.
    time.sleep(0.01)
    if request.param == 'rust':
        server = subprocess.Popen(["target/debug/litemq"])
    else:
        server = subprocess.Popen(["build/litemq"], env={"PORT": "42099"})
    # Sleep here to give the server time to start before the tests run.
    time.sleep(0.01)
    try:
        yield request.param
    finally:
        server.kill()


def new_client(param):
    if param == 'rust':
        channel = Channel('127.0.0.1', 42069)
    else:
        channel = Channel('127.0.0.1', 42099)
    return LiteMQ(channel)
