package database

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Health() map[string]string
}

type service struct {
	db *redis.Client
}

func New() Service {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	address  := os.Getenv("CACHEFLOW_REDIS_ADDRESS")
	port     := os.Getenv("CACHEFLOW_REDIS_PORT")
	password := os.Getenv("CACHEFLOW_REDIS_PASSWORD")
	// database := os.Getenv("CACHEFLOW_REDIS_DATABASE")
	database, check := os.LookupEnv("CACHEFLOW_REDIS_DATABASE")
	if !check {
		log.Fatal("database env not present")
	}
	num, err := strconv.Atoi(database)
	if err != nil {
		log.Fatalf("database incorrect %v", err)
	}

	fullAddress := fmt.Sprintf("%s:%s", address, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
		DB:       num,
	})

	s := &service{
		db: rdb,
	}

	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := make(map[string]string)

	stats = s.checkRedisHealth(ctx, stats)

	return stats
}

func (s *service) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	pong, err := s.db.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	stats["redis_status"] = "up"
	stats["redis_message"] = "It's Healthy"
	stats["redis_ping_response"] = pong

	info, err := s.db.Info(ctx).Result()
	if err != nil {
		stats["redis_message"] = fmt.Sprintf("Failed to retrieve Redis info: %v", err)
		return stats
	}

	redisInfo := parseRedisInfo(info)

	poolStats := s.db.PoolStats()

	stats["redis_version"] = redisInfo["redis_version"]
	stats["redis_mode"] = redisInfo["redis_mode"]
	stats["redis_connected_clients"] = redisInfo["connected_clients"]
	stats["redis_used_memory"] = redisInfo["used_memory"]
	stats["redis_used_memory_peak"] = redisInfo["used_memory_peak"]
	stats["redis_uptime_in_seconds"] = redisInfo["uptime_in_seconds"]
	stats["redis_hits_connections"] = strconv.FormatUint(uint64(poolStats.Hits), 10)
	stats["redis_misses_connections"] = strconv.FormatUint(uint64(poolStats.Misses), 10)
	stats["redis_timeouts_connections"] = strconv.FormatUint(uint64(poolStats.Timeouts), 10)
	stats["redis_total_connections"] = strconv.FormatUint(uint64(poolStats.TotalConns), 10)
	stats["redis_idle_connections"] = strconv.FormatUint(uint64(poolStats.IdleConns), 10)
	stats["redis_stale_connections"] = strconv.FormatUint(uint64(poolStats.StaleConns), 10)
	stats["redis_max_memory"] = redisInfo["maxmemory"]

	activeConns := uint64(math.Max(float64(poolStats.TotalConns-poolStats.IdleConns), 0))
	stats["redis_active_connections"] = strconv.FormatUint(activeConns, 10)

	poolSize := s.db.Options().PoolSize
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	poolSizePercentage := float64(connectedClients) / float64(poolSize) * 100
	stats["redis_pool_size_percentage"] = fmt.Sprintf("%.2f%%", poolSizePercentage)

	return s.evaluateRedisStats(redisInfo, stats)
}

func (s *service) evaluateRedisStats(redisInfo, stats map[string]string) map[string]string {
	poolSize := s.db.Options().PoolSize
	poolStats := s.db.PoolStats()
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	highConnectionThreshold := int(float64(poolSize) * 0.8)

	// Check if the number of connected clients is high.
	if connectedClients > highConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of connected clients"
	}

	// Check if the number of stale connections exceeds a threshold.
	minStaleConnectionsThreshold := 500
	if int(poolStats.StaleConns) > minStaleConnectionsThreshold {
		stats["redis_message"] = fmt.Sprintf("Redis has %d stale connections.", poolStats.StaleConns)
	}

	// Check if Redis is using a significant amount of memory.
	usedMemory, _ := strconv.ParseInt(redisInfo["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(redisInfo["maxmemory"], 10, 64)
	if maxMemory > 0 {
		usedMemoryPercentage := float64(usedMemory) / float64(maxMemory) * 100
		if usedMemoryPercentage >= 90 {
			stats["redis_message"] = "Redis is using a significant amount of memory"
		}
	}

	// Check if Redis has been recently restarted.
	uptimeInSeconds, _ := strconv.ParseInt(redisInfo["uptime_in_seconds"], 10, 64)
	if uptimeInSeconds < 3600 {
		stats["redis_message"] = "Redis has been recently restarted"
	}

	// Check if the number of idle connections is high.
	idleConns := int(poolStats.IdleConns)
	highIdleConnectionThreshold := int(float64(poolSize) * 0.7)
	if idleConns > highIdleConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of idle connections"
	}

	// Check if the connection pool utilization is high.
	poolUtilization := float64(poolStats.TotalConns-poolStats.IdleConns) / float64(poolSize) * 100
	highPoolUtilizationThreshold := 90.0
	if poolUtilization > highPoolUtilizationThreshold {
		stats["redis_message"] = "Redis connection pool utilization is high"
	}

	return stats
}

// parseRedisInfo parses the Redis info response and returns a map of key-value pairs.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

