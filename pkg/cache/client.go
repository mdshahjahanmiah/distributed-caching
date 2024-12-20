package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"hash/fnv"
	"sync"
)

var (
	ctx        = context.Background()
	redisNodes []*redis.Client
	mutex      sync.Mutex
)

func InitializeNodes() {
	mutex.Lock()
	defer mutex.Unlock()

	redisNodes = []*redis.Client{
		newRedisClient("localhost:6379"), // North America
		newRedisClient("localhost:6380"), // Asia
		newRedisClient("localhost:6381"), // Europe
	}
}

func GetClientForKey(key string) *redis.Client {
	hash := computeHash(key)
	index := hash % uint32(len(redisNodes))
	return redisNodes[index]
}

func computeHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func newRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func Context() context.Context {
	return ctx
}

func Cleanup() {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range redisNodes {
		client.Close()
	}
}
