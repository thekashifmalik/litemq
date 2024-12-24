import time
import subprocess
import asyncio
import pytest
import pytest_asyncio
import nest_asyncio
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


async def test_concurrent(server):
    client = new_client(server)

    async def task():
        assert await client.enqueue('test-blocking', b'message') == 0

    # This test is flakey, as are all tests for concurrency. We run it a couple of time since it does reliably catch
    # errors.
    for _ in range(100):
        asyncio.create_task(task())
        assert await client.dequeue('test-blocking') == b'message'


async def test_performance_length(server, aio_benchmark):
    client = new_client(server)
    aio_benchmark(client.length, 'test-benchmark')


async def test_performance_enqueue(server, aio_benchmark):
    client = new_client(server)
    aio_benchmark(client.enqueue, 'test-benchmark', b'message')


async def test_performance_enqueue_and_dequeue(server, aio_benchmark):
    client = new_client(server)
    async def task():
        await client.enqueue('test-benchmark', b'message')
        await client.dequeue('test-benchmark')
    aio_benchmark(task)


TEST_PORT = 42099
TEST_ENV = {
    'PORT': str(TEST_PORT),
}


@pytest.fixture(params=['go', 'rust'])
def server(request):
    # Sleep here so the server has enough time to give up the port between tests.
    time.sleep(0.01)
    if request.param == 'rust':
        server = subprocess.Popen(["target/release/litemq"], env=TEST_ENV)
    else:
        server = subprocess.Popen(["build/litemq"], env=TEST_ENV)
    # Sleep here to give the server time to start before the tests run.
    time.sleep(0.02)
    try:
        yield request.param
    finally:
        server.kill()


def new_client(param):
    if param == 'rust':
        channel = Channel('127.0.0.1', TEST_PORT)
    else:
        channel = Channel('127.0.0.1', TEST_PORT)

    return LiteMQ(channel)


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
