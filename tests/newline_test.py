import asyncio
import pytest

from .core import server
from .core import flushed_client


async def test_messages_with_newlines(server):
    """Test that messages containing newlines are handled correctly"""
    client = await flushed_client()

    # Test messages with various newline combinations
    test_messages = [
        b"Hello\nWorld",                    # Single newline
        b"Line1\nLine2\nLine3",            # Multiple newlines
        b"\nStarts with newline",          # Starts with newline
        b"Ends with newline\n",            # Ends with newline
        b"Normal message without newlines", # Control message
        b"",                               # Empty message
    ]

    # Enqueue all test messages
    queue_name = 'newline_test'
    for i, msg in enumerate(test_messages):
        length = await client.enqueue(queue_name, msg)
        assert length == i + 1

    # Verify queue length
    final_length = await client.length(queue_name)
    assert final_length == len(test_messages)

    # Dequeue and verify all messages
    for expected_msg in test_messages:
        dequeued = await client.dequeue(queue_name)
        assert dequeued == expected_msg

    # Verify queue is now empty
    final_length = await client.length(queue_name)
    assert final_length == 0


async def test_message_with_only_newlines(server):
    """Test message that contains only newlines"""
    client = await flushed_client()

    message = b"\n\n\n"
    queue_name = 'only_newlines_test'

    assert await client.enqueue(queue_name, message) == 1
    assert await client.length(queue_name) == 1
    assert await client.dequeue(queue_name) == message
    assert await client.length(queue_name) == 0
