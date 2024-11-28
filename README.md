# LiteMQ
Free & open source, high-performance, distributed message queue system, generic in nature, but intended for use in
background task processing for dynamic web applications.

## Why
Redis [went rogue](https://www.reddit.com/r/redis/comments/1bjs7bo/redis_is_switching_away_from_opensource_licensing/)
so I needed to replace it. [Memcached](https://memcached.org/) works great as a cache and key/value store but it does
not provide an equivalent to Redis lists, which are often used as message queues for background task processing.
