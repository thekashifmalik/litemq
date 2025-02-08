import time
import asyncio

from .core import server, flushed_client



async def test_enqueue_while_blocking_dequeue(server):
    client = await flushed_client(server)

    async def task():
        time.sleep(0.01)
        assert await client.enqueue('test-blocking', b'message') == 0

    # This test is flakey, as are all tests for concurrency. We run it a couple of time since it does reliably catch
    # errors.
    for _ in range(100):
        asyncio.create_task(task())
        assert await client.dequeue('test-blocking') == b'message'
