package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func testCluster(log *zap.Logger, data []Entry) {
	log.Sugar().Info("Trying connect redis cluster...")
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:6373", "localhost:6374", "localhost:6375"},
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Sugar().Error("Cannot connect redis cluster")
		return
	}

	defer func(rdb *redis.ClusterClient) {
		if err := rdb.Close(); err != nil {
			log.Sugar().Error(err)
		}
	}(rdb)
	log.Sugar().Info("Successfully connected to redis cluster")

	var spentMillis int64

	// HSET
	{
		p := rdb.TxPipeline()
		for _, entry := range data {
			id := "photos:" + entry.PhotoID
			p.HSet(ctx, id, "business_id", entry.BusinessID)
			p.HSet(ctx, id, "caption", entry.Caption)
			p.HSet(ctx, id, "label", entry.Label)
		}
		start := time.Now()
		_, _ = p.Exec(ctx)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("HSET time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// HGET
	{
		p := rdb.TxPipeline()
		for _, entry := range data {
			id := "photos:" + entry.PhotoID
			_ = p.HGetAll(ctx, id)
		}
		start := time.Now()
		_, _ = p.Exec(ctx)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("HGET time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// ZADD
	{
		p := rdb.TxPipeline()
		for _, entry := range data {
			bid := entry.BusinessID
			_ = p.ZAdd(ctx, "businesses", redis.Z{Member: bid, Score: 1})
		}
		start := time.Now()
		_, _ = p.Exec(ctx)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("ZADD time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// ZSCORE
	{
		p := rdb.TxPipeline()
		for _, entry := range data {
			bid := entry.BusinessID
			_ = p.Do(ctx, "ZSCORE", "businesses", bid)
		}
		start := time.Now()
		_, _ = p.Exec(ctx)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("ZSCORE time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// LPUSH
	{
		p := rdb.TxPipeline()
		for _, entry := range data {
			_ = p.LPush(ctx, "photos", entry)
		}
		start := time.Now()
		_, _ = p.Exec(ctx)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("LPUSH time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// LPOP
	{
		p := rdb.TxPipeline()
		for range data {
			_ = p.LPop(ctx, "photos")
		}
		start := time.Now()
		_, _ = p.Exec(ctx)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("LPOP time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// String (SET)
	raw, _ := json.Marshal(map[string]any{"data": data})
	{
		start := time.Now()
		rdb.Set(ctx, "raw_json", raw, 0)
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("SET time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	// String (GET)
	{
		start := time.Now()
		rdb.Get(ctx, "raw_json")
		spentMillis = time.Since(start).Milliseconds()
	}
	log.Sugar().Infof("GET time: %.1f s", float64(spentMillis)/1000)
	_ = log.Sync()

	_ = rdb.FlushAll(ctx)
}
