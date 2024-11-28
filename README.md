# LiteMQ
Free & open source, high-performance, distributed message queue system, generic in nature, but intended for use in
background task processing for dynamic web applications.


## Why
Redis [went rogue](https://www.reddit.com/r/redis/comments/1bjs7bo/redis_is_switching_away_from_opensource_licensing/)
so I needed somethign to replace it. Memcached works great as a cache and key/value store but it does not replace other
data-structures that Redis provided. This software replaces usage of Redis lists as message queues, something that's
often used for background task processing.
