# litemq
Free & open source in-memory message queue, generic in nature, but intended for use in background task processing for
web applications.

![Screenshot](./screenshot.png)

## Quickstart
Fastest way to use this software is to run it using docker compose:

```yaml
services:

  litemq:
    image: thekashifmalik/litemq:latest
    ports:
      - 42090:42090

```

You can also use the docker image directly:

```bash
docker run --rm -it -p 42090:42090 thekashifmalik/litemq:latest
```

## Config
You can set the following environment variables:

- `PORT`: Change the port the server uses.
- `LOG_LEVEL`: Set to `debug` to enable debug logging.


## Why
I needed a message queue for background task processing. [Redis](https://redis.io/) is often used for this purpose, by
way of Redis lists, but that project is [caught up](https://www.reddit.com/r/redis/comments/1bjs7bo/redis_is_switching_away_from_opensource_licensing/)
in licensing issues. [ValKey](https://valkey.io/) seems like a decent alternative and [RabbitMQ](https://www.rabbitmq.com/)
is the traditional message broker, but I wanted something lighter and more focused, similar to [Memcached](https://memcached.org/),
which works great as a cache and not much else.

## References
- [Website](https://litemq.com/)
- [Repository](https://github.com/thekashifmalik/litemq)
