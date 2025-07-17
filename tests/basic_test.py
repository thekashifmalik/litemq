import asyncio
import pytest

from .core import server
from .core import flushed_client



async def test_server(server):
    assert True


async def test_health(server):
    client = await flushed_client()
    await client.health()


async def test_enqueue_single(server):
    client = await flushed_client()
    assert await client.enqueue('test', b'message-1') == 1


async def test_enqueue_multiple(server):
    client = await flushed_client()
    assert await client.enqueue('test', b'message-1') == 1
    assert await client.enqueue('test', b'message-2') == 2


async def test_enqueue_different_queues(server):
    client = await flushed_client()
    assert await client.enqueue('test-1', b'message-1') == 1
    assert await client.enqueue('test-2', b'message-1') == 1


async def test_enqueue_different_queues_multiple(server):
    client = await flushed_client()
    assert await client.enqueue('test-1', b'message-1') == 1
    assert await client.enqueue('test-1', b'message-2') == 2
    assert await client.enqueue('test-2', b'message-1') == 1
    assert await client.enqueue('test-2', b'message-2') == 2


async def test_enqueue_large(server):
    client = await flushed_client()
    assert await client.enqueue('test', b'message' * 1000) == 1


async def test_length(server):
    client = await flushed_client()
    assert await client.length('test') == 0
    await client.enqueue('test', b'message')
    assert await client.length('test') == 1


async def test_dequeue_single(server):
    client = await flushed_client()
    await client.enqueue('test', b'message')
    assert await client.dequeue('test') == b'message'


async def test_dequeue_multiple(server):
    client = await flushed_client()
    await client.enqueue('test', b'message-1')
    await client.enqueue('test', b'message-2')
    assert await client.dequeue('test') == b'message-1'
    assert await client.dequeue('test') == b'message-2'


async def test_dequeue_different_queues(server):
    client = await flushed_client()
    await client.enqueue('test-1', b'message-1')
    await client.enqueue('test-2', b'message-2')
    assert await client.dequeue('test-1') == b'message-1'
    assert await client.dequeue('test-2') == b'message-2'


async def test_dequeue_different_queues_multiple(server):
    client = await flushed_client()
    await client.enqueue('test-1', b'message-1')
    await client.enqueue('test-1', b'message-2')
    await client.enqueue('test-2', b'message-1')
    await client.enqueue('test-2', b'message-2')
    assert await client.dequeue('test-1') == b'message-1'
    assert await client.dequeue('test-1') == b'message-2'
    assert await client.dequeue('test-2') == b'message-1'
    assert await client.dequeue('test-2') == b'message-2'


async def test_dequeue_blocks(server):
    client = await flushed_client()
    with pytest.raises(asyncio.TimeoutError):
        await asyncio.wait_for(client.dequeue('test-1'), timeout=0.1)


async def test_purge(server):
    client = await flushed_client()
    await client.enqueue('test', b'message')
    assert await client.length('test') == 1
    assert await client.purge('test') == 1
    assert await client.length('test') == 0


async def test_flush(server):
    client = await flushed_client()
    await client.enqueue('test-1', b'message')
    await client.enqueue('test-2', b'message')
    assert await client.length('test-1') == 1
    assert await client.length('test-2') == 1
    assert await client.flush()
    assert await client.length('test-1') == 0
    assert await client.length('test-2') == 0
