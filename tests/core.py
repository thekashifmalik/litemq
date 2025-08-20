from contextlib import contextmanager
import time
import subprocess

import pytest
from grpclib.client import Channel

from litemq.client import LiteMQ



TEST_PORT = 42099
TEST_ENV = {
    'PORT': str(TEST_PORT),
}

@contextmanager
def running_server():
    server = subprocess.Popen(["target/release/litemq"], env=TEST_ENV)
    # Sleep here to give the server time to start before the tests run.
    time.sleep(0.02)
    # Check if the server started successfully
    if server.poll() is not None:
        raise RuntimeError("Server failed to start")
    try:
        yield
    finally:
        server.kill()


@pytest.fixture
def server(request):
    # Sleep here so the server has enough time to give up the port between tests.
    time.sleep(0.01)
    with running_server():
        yield


def new_client():
    channel = Channel('127.0.0.1', TEST_PORT)
    return LiteMQ(channel)


async def flushed_client():
    client = new_client()
    await client.flush()
    return client
