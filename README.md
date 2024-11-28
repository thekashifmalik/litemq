# litemq
Simple and straightforward message queue.

## Why
Redis went rogue so I needed somethign to replace it's usage in my personal projects. Memcached works great as a cache
and key/value store but it does not replace other data-structures. This software will replace my usage of Redis lists as
message queues, somethibng I used for background tasks in Django.
