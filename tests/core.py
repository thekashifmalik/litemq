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
def running_server(param):
    if param == 'rust':
        server = subprocess.Popen(["target/release/litemq"], env=TEST_ENV)
    else:
        server = subprocess.Popen(["build/litemq"], env=TEST_ENV)
    # Sleep here to give the server time to start before the tests run.
    time.sleep(0.02)
    try:
        yield
    finally:
        server.kill()


@pytest.fixture(params=['go', 'rust'])
def server(request):
    # Sleep here so the server has enough time to give up the port between tests.
    time.sleep(0.01)
    with running_server(request.param):
        yield request.param


def new_client(param):
    channel = Channel('127.0.0.1', TEST_PORT)
    return LiteMQ(channel)
