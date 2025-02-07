from .core import new_client
from .core import running_server


async def test_persistence_go():
    await run_persistence_check('go')


async def test_persistence_rust():
    await run_persistence_check('rust')


async def run_persistence_check(param):
    with running_server(param):
        client = new_client(param)
        assert await client.enqueue('test-persistence', b'message') == 1
        assert await client.length('test-persistence') == 1
        client.channel.close()

    with running_server(param):
        client = new_client(param)
        assert await client.length('test-persistence') == 1
        assert await client.dequeue('test-persistence') == b'message'
