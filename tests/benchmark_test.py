import asyncio
import pytest_asyncio
import nest_asyncio
from .core import server, flushed_client


async def test_length(server, aio_benchmark):
    client = await flushed_client()
    aio_benchmark(client.length, 'test-benchmark-length')


async def test_enqueue(server, aio_benchmark):
    client = await flushed_client()
    aio_benchmark(client.enqueue, 'test-benchmark-enqueue', b'message')


async def test_enqueue_large(server, aio_benchmark):
    client = await flushed_client()
    aio_benchmark(client.enqueue, 'test-benchmark-enqueue-large', b'message' * 1000)


async def test_enqueue_and_dequeue(server, aio_benchmark):
    client = await flushed_client()
    async def task():
        await client.enqueue('test-benchmark-enqueue-and-dequeue', b'message')
        await client.dequeue('test-benchmark-enqueue-and-dequeue')
    aio_benchmark(task)


async def test_dequeues_on_large_queue(server, aio_benchmark):
    client = await flushed_client()
    for _ in range(10000000):
        client.enqueue('test-benchmark-dequeues-on-large-queue', b'message' * 10000)

    async def task():
        await client.enqueue('test-benchmark-dequeues-on-large-queue', b'message')
        await client.dequeue('test-benchmark-dequeues-on-large-queue')
    aio_benchmark(task)


# COPYPASTA: https://github.com/ionelmc/pytest-benchmark/issues/66#issuecomment-1137005280
# This is a workaround for the pytest-benchmark not supporting async functions.
@pytest_asyncio.fixture
async def aio_benchmark(benchmark):
    nest_asyncio.apply()

    def _wrapper(func, *args, **kwargs):
        if asyncio.iscoroutinefunction(func):
            @benchmark
            def _():
                return asyncio.run(func(*args, **kwargs))
                # return event_loop.run_until_complete(func(*args, **kwargs))
        else:
            benchmark(func, *args, **kwargs)

    return _wrapper
