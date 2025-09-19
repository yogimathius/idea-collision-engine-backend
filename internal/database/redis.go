package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"idea-collision-engine-api/internal/models"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Cache keys
const (
	KeyCollisionDomains = "collision:domains:%s"        // collision:domains:tier
	KeyUserUsage        = "user:usage:%s"               // user:usage:user_id
	KeyCollisionResult  = "collision:result:%s"         // collision:result:hash
	KeyRateLimit        = "rate:limit:%s:%d"            // rate:limit:user_id:window
)

// Cache collision domains by tier
func (r *RedisClient) CacheCollisionDomains(tier string, domains []models.CollisionDomain, expiration time.Duration) error {
	key := fmt.Sprintf(KeyCollisionDomains, tier)
	
	data, err := json.Marshal(domains)
	if err != nil {
		return fmt.Errorf("failed to marshal collision domains: %w", err)
	}
	
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisClient) GetCachedCollisionDomains(tier string) ([]models.CollisionDomain, error) {
	key := fmt.Sprintf(KeyCollisionDomains, tier)
	
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var domains []models.CollisionDomain
	err = json.Unmarshal([]byte(data), &domains)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal collision domains: %w", err)
	}
	
	return domains, nil
}

// Cache user usage for rate limiting
func (r *RedisClient) CacheUserUsage(userID string, usage *models.UserUsage, expiration time.Duration) error {
	key := fmt.Sprintf(KeyUserUsage, userID)
	
	data, err := json.Marshal(usage)
	if err != nil {
		return fmt.Errorf("failed to marshal user usage: %w", err)
	}
	
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisClient) GetCachedUserUsage(userID string) (*models.UserUsage, error) {
	key := fmt.Sprintf(KeyUserUsage, userID)
	
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var usage models.UserUsage
	err = json.Unmarshal([]byte(data), &usage)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user usage: %w", err)
	}
	
	return &usage, nil
}

// Cache collision results for similar requests
func (r *RedisClient) CacheCollisionResult(inputHash string, result *models.CollisionResult, expiration time.Duration) error {
	key := fmt.Sprintf(KeyCollisionResult, inputHash)
	
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal collision result: %w", err)
	}
	
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisClient) GetCachedCollisionResult(inputHash string) (*models.CollisionResult, error) {
	key := fmt.Sprintf(KeyCollisionResult, inputHash)
	
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var result models.CollisionResult
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal collision result: %w", err)
	}
	
	return &result, nil
}

// Rate limiting using sliding window
func (r *RedisClient) CheckRateLimit(userID string, windowSeconds int, limit int) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-time.Duration(windowSeconds) * time.Second)
	
	key := fmt.Sprintf(KeyRateLimit, userID, windowSeconds)
	
	// Remove old entries
	err := r.client.ZRemRangeByScore(r.ctx, key, "0", fmt.Sprintf("%d", windowStart.Unix())).Err()
	if err != nil {
		return false, err
	}
	
	// Count current requests
	count, err := r.client.ZCard(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	if int(count) >= limit {
		return false, nil // Rate limit exceeded
	}
	
	// Add current request
	err = r.client.ZAdd(r.ctx, key, redis.Z{
		Score:  float64(now.Unix()),
		Member: now.UnixNano(),
	}).Err()
	
	if err != nil {
		return false, err
	}
	
	// Set expiration on the key
	r.client.Expire(r.ctx, key, time.Duration(windowSeconds)*time.Second)
	
	return true, nil
}

// Get rate limit status
func (r *RedisClient) GetRateLimitStatus(userID string, windowSeconds int, limit int) (int, time.Duration, error) {
	now := time.Now()
	windowStart := now.Add(-time.Duration(windowSeconds) * time.Second)
	
	key := fmt.Sprintf(KeyRateLimit, userID, windowSeconds)
	
	// Remove old entries
	err := r.client.ZRemRangeByScore(r.ctx, key, "0", fmt.Sprintf("%d", windowStart.Unix())).Err()
	if err != nil {
		return 0, 0, err
	}
	
	// Count current requests
	count, err := r.client.ZCard(r.ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}
	
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	
	// Get time until oldest request expires
	var resetTime time.Duration
	if count > 0 {
		oldest, err := r.client.ZRange(r.ctx, key, 0, 0).Result()
		if err == nil && len(oldest) > 0 {
			oldestTime := time.Unix(0, 0) // fallback
			if member, err := r.client.ZScore(r.ctx, key, oldest[0]).Result(); err == nil {
				oldestTime = time.Unix(int64(member), 0)
			}
			resetTime = oldestTime.Add(time.Duration(windowSeconds) * time.Second).Sub(now)
			if resetTime < 0 {
				resetTime = 0
			}
		}
	}
	
	return remaining, resetTime, nil
}

// Invalidate cache entries
func (r *RedisClient) InvalidateCollisionDomains(tier string) error {
	key := fmt.Sprintf(KeyCollisionDomains, tier)
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) InvalidateUserUsage(userID string) error {
	key := fmt.Sprintf(KeyUserUsage, userID)
	return r.client.Del(r.ctx, key).Err()
}

// Health check
func (r *RedisClient) Ping() error {
	return r.client.Ping(r.ctx).Err()
}