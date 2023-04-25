# Redis & Redis Cluster performance

## Environment

- `docker-compose.yml`: redis and redis-cluster services
- `docker/`: cluster nodes configurations
- `yelp_photos.json`: open YELP dataset
- `prof`: Golang implementation

## How To Run

```shell
docker-compose up redis-single redis-cluster
go build prof
./prof.exe
```

## Golang Implementations Details

- `HSET`/`HGET`, `ZADD`/`ZSCORE`, `LPUSH`/`LPOP`, `SET`/`GET` operations
- pipelining for 1 RTT per whole dataset operation
- `zap` logging

## Raw Logs

- System: Win10 amd64

```
2023-04-24T22:13:28.230+0300    INFO    prof/main.go:40 Successfully read test data
2023-04-24T22:13:28.269+0300    INFO    prof/single.go:13       Trying connect redis...
2023-04-24T22:13:28.285+0300    INFO    prof/single.go:26       Successfully connected to redis
2023-04-24T22:13:40.733+0300    INFO    prof/single.go:43       HSET time: 12.2 s
2023-04-24T22:13:41.162+0300    INFO    prof/single.go:57       HGET time: 0.3 s
2023-04-24T22:13:41.843+0300    INFO    prof/single.go:71       ZADD time: 0.6 s
2023-04-24T22:13:42.386+0300    INFO    prof/single.go:85       ZSCORE time: 0.5 s
2023-04-24T22:13:45.678+0300    INFO    prof/single.go:98       LPUSH time: 3.1 s
2023-04-24T22:13:46.022+0300    INFO    prof/single.go:111      LPOP time: 0.3 s
2023-04-24T22:13:46.390+0300    INFO    prof/single.go:121      SET time: 0.3 s
2023-04-24T22:13:47.038+0300    INFO    prof/single.go:130      GET time: 0.6 s
2023-04-24T22:13:47.079+0300    INFO    prof/cluster.go:13      Trying connect redis cluster...
2023-04-24T22:13:47.093+0300    INFO    prof/cluster.go:28      Successfully connected to redis cluster
2023-04-24T22:14:10.158+0300    INFO    prof/cluster.go:45      HSET time: 22.7 s
2023-04-24T22:14:29.891+0300    INFO    prof/cluster.go:59      HGET time: 19.6 s
2023-04-24T22:14:33.337+0300    INFO    prof/cluster.go:73      ZADD time: 3.3 s
2023-04-24T22:14:36.692+0300    INFO    prof/cluster.go:87      ZSCORE time: 3.3 s
2023-04-24T22:14:37.283+0300    INFO    prof/cluster.go:100     LPUSH time: 0.4 s
2023-04-24T22:14:37.655+0300    INFO    prof/cluster.go:113     LPOP time: 0.3 s
2023-04-24T22:14:39.130+0300    INFO    prof/cluster.go:123     SET time: 1.3 s
2023-04-24T22:14:39.293+0300    INFO    prof/cluster.go:132     GET time: 0.2 s
```

## Summary

- `GET` requests can be balanced between nodes (x3 better performance for cluster)
- `HSET` / `HGET` slowdown >10 for cluster
- `ZADD` / `ZSCORE` x5 slowdown
- Faster `LPUSH` and the same `LPOP` with cluster
