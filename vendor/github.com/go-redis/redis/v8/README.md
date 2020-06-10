# Redis client for Golang

[![Build Status](https://travis-ci.org/go-redis/redis.png?branch=master)](https://travis-ci.org/go-redis/redis)
[![GoDoc](https://godoc.org/github.com/go-redis/redis?status.svg)](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc)
[![Airbrake](https://img.shields.io/badge/kudos-airbrake.io-orange.svg)](https://airbrake.io)

Supports:

- Redis 3 commands except QUIT, MONITOR, SLOWLOG and SYNC.
- Automatic connection pooling with [circuit breaker](https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern) support.
- [Pub/Sub](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#PubSub).
- [Transactions](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#example-Client-TxPipeline).
- [Pipeline](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#example-Client-Pipeline) and [TxPipeline](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#example-Client-TxPipeline).
- [Scripting](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#Script).
- [Timeouts](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#Options).
- [Redis Sentinel](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#NewFailoverClient).
- [Redis Cluster](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#NewClusterClient).
- [Cluster of Redis Servers](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#example-NewClusterClient--ManualSetup) without using cluster mode and Redis Sentinel.
- [Ring](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#NewRing).
- [Instrumentation](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#ex-package--Instrumentation).
- [Cache friendly](https://github.com/go-redis/cache).
- [Rate limiting](https://github.com/go-redis/redis_rate).
- [Distributed Locks](https://github.com/bsm/redislock).

API docs: https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc.
Examples: https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#pkg-examples.

## Installation

go-redis requires a Go version with [Modules](https://github.com/golang/go/wiki/Modules) support and uses import versioning. So please make sure to initialize a Go module before installing go-redis:

```shell
go mod init github.com/my/repo
go get github.com/go-redis/redis/v8
```

Import:

```go
import "github.com/go-redis/redis/v8"
```

## Quickstart

```go
func ExampleNewClient() {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := client.Ping(ctx).Result()
    fmt.Println(pong, err)
    // Output: PONG <nil>
}

func ExampleClient() {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    err := client.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := client.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := client.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
    // Output: key value
    // key2 does not exist
}
```

## Howto

Please go through [examples](https://pkg.go.dev/github.com/go-redis/redis/v8?tab=doc#pkg-examples) to get an idea how to use this package.

## Look and feel

Some corner cases:

```go
// SET key value EX 10 NX
set, err := client.SetNX(ctx, "key", "value", 10*time.Second).Result()

// SORT list LIMIT 0 2 ASC
vals, err := client.Sort(ctx, "list", &redis.Sort{Offset: 0, Count: 2, Order: "ASC"}).Result()

// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
vals, err := client.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
    Min: "-inf",
    Max: "+inf",
    Offset: 0,
    Count: 2,
}).Result()

// ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
vals, err := client.ZInterStore(ctx, "out", &redis.ZStore{
    Keys: []string{"zset1", "zset2"},
    Weights: []int64{2, 3}
}).Result()

// EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
vals, err := client.Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()

// custom command
res, err := client.Do(ctx, "set", "key", "value").Result()
```

## See also

- [Golang PostgreSQL ORM](https://github.com/go-pg/pg)
- [Golang msgpack](https://github.com/vmihailenco/msgpack)
- [Golang message task queue](https://github.com/vmihailenco/taskq)
