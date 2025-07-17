from .core import new_client
from .core import running_server


async def test_persistence():
    with running_server():
        client = new_client()
        await client.flush()
        assert await client.enqueue('test-persistence', b'message') == 1
        assert await client.length('test-persistence') == 1
        client.channel.close()

    with running_server():
        client = new_client()
        assert await client.length('test-persistence') == 1
        assert await client.dequeue('test-persistence') == b'message'
        await client.flush()
