package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"hash/fnv"
	"sort"
	"sync"
	"time"
)

// Global variables for managing the Redis cluster and cache
var (
	ctx        = context.Background()
	hashRing   *ConsistentHash              // Consistent hashing ring for node management
	redisNodes map[string]*redis.Client     // Map of Redis nodes by region
	mutex      sync.Mutex                   // Mutex for thread-safe operations
	cacheItems = make(map[string]CacheItem) // Simulated cache for tracking key access patterns
)

// CacheItem represents an item in the cache with metadata
type CacheItem struct {
	Value    string    // Stored value (optional, not currently used)
	LastUsed time.Time // Timestamp of the last access
}

// ConsistentHash manages consistent hashing for Redis node allocation
type ConsistentHash struct {
	replicas     int               // Number of replicas for each node in the hash ring
	hashMap      map[uint32]string // Mapping of hash values to nodes
	sortedHashes []uint32          // Sorted list of hash values for efficient lookups
}

// NewConsistentHash initializes a new consistent hash ring
func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		replicas:     replicas,
		hashMap:      make(map[uint32]string),
		sortedHashes: make([]uint32, 0),
	}
}

// AddNode adds a new node to the hash ring with the specified number of replicas
func (ch *ConsistentHash) AddNode(node string) {
	for i := 0; i < ch.replicas; i++ {
		hash := ch.hashKey(fmt.Sprintf("%s:%d", node, i))
		ch.hashMap[hash] = node
		ch.sortedHashes = append(ch.sortedHashes, hash)
	}
	sort.Slice(ch.sortedHashes, func(i, j int) bool { return ch.sortedHashes[i] < ch.sortedHashes[j] })
}

// GetNode retrieves the node responsible for the given key
func (ch *ConsistentHash) GetNode(key string) string {
	if len(ch.hashMap) == 0 {
		return ""
	}

	hash := ch.hashKey(key)
	idx := sort.Search(len(ch.sortedHashes), func(i int) bool {
		return ch.sortedHashes[i] >= hash
	})

	if idx == len(ch.sortedHashes) {
		idx = 0
	}

	return ch.hashMap[ch.sortedHashes[idx]]
}

// hashKey generates a hash value for the given key
func (ch *ConsistentHash) hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// InitializeNodes initializes Redis clients and consistent hashing
func InitializeNodes() {
	mutex.Lock()
	defer mutex.Unlock()

	// Set up Redis clients for each region
	redisNodes = map[string]*redis.Client{
		"Europe":       newRedisClient("localhost:6379"),
		"Asia":         newRedisClient("localhost:6380"),
		"NorthAmerica": newRedisClient("localhost:6381"),
	}

	// Initialize consistent hash ring with 100 replicas per node
	hashRing = NewConsistentHash(100)
	for node := range redisNodes {
		hashRing.AddNode(node)
	}
}

// GetClientForKey selects a Redis client based on the key and user's region
func GetClientForKey(key string, userRegion string) *redis.Client {
	// Use consistent hashing to get the node for the key
	node := hashRing.GetNodeAndUpdateLastUse(key)

	// Prefer a node in the user's region if available
	if client, exists := redisNodes[userRegion]; exists {
		if node == userRegion {
			return client
		}
		// Look for a node explicitly in the user's region
		for k, v := range redisNodes {
			if k == userRegion {
				return v
			}
		}
	}
	// Default to the client from consistent hashing if no region-specific match found
	return redisNodes[node]
}

// GetNodeAndUpdateLastUse retrieves a node and updates the cache access time
func (ch *ConsistentHash) GetNodeAndUpdateLastUse(key string) string {
	node := ch.GetNode(key)
	if item, exists := cacheItems[key]; exists {
		item.LastUsed = time.Now()
		cacheItems[key] = item
	} else {
		cacheItems[key] = CacheItem{
			LastUsed: time.Now(),
		}
	}
	return node
}

// newRedisClient creates a new Redis client with default settings
func newRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // No password set for these examples
		DB:       0,  // Default database
	})
}

// Context returns the global context instance
func Context() context.Context {
	return ctx
}

// Cleanup closes all Redis clients to free resources
func Cleanup() {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range redisNodes {
		client.Close()
	}
}
