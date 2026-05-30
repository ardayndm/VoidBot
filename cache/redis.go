package cache

import (
	config "VoidBot/config/core"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// Redis - Redis veritabanı işlemlerini yönetir
var Redis *RedisCache

// NewRedis yeni bir Redis bağlantısı oluşturur
// Kullanım: cache.NewRedis("localhost:6379", "", 0)
func NewRedis(cfg config.RedisConfig) error {
	addr := cfg.Host + ":" + cfg.Port

	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Bağlantıyı test et
	if err := cli.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis bağlantı hatası: %w", err)
	}

	Redis = &RedisCache{
		client: cli,
		ctx:    ctx,
	}
	return nil
}

// Get - anahtara ait değeri getirir, anahtar yoksa "", nil döner
func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// Set - anahtara değer atar, expTime=0 ise süresiz kalır
func (r *RedisCache) Set(key string, value interface{}, expTime time.Duration) (bool, error) {
	err := r.client.Set(r.ctx, key, value, expTime).Err()
	return err == nil, err
}

// Delete - anahtarı siler
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists - anahtar var mı kontrol eder
func (r *RedisCache) Exists(key string) (bool, error) {
	n, err := r.client.Exists(r.ctx, key).Result()
	return n > 0, err
}

// SetIfNotExists - anahtar yoksa set eder, varsa false döner (cooldown için ideal)
func (r *RedisCache) SetIfNotExists(key string, value interface{}, expTime time.Duration) (bool, error) {
	return r.client.SetNX(r.ctx, key, value, expTime).Result()
}

// Incr - anahtarın değerini 1 arttırır
func (r *RedisCache) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Decr - anahtarın değerini 1 azaltır
func (r *RedisCache) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// Expire - var olan anahtarın süresini değiştirir
func (r *RedisCache) ExpireSet(key string, expTime time.Duration) error {
	return r.client.Expire(r.ctx, key, expTime).Err()
}

// Close - bağlantıyı kapatır
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Ping - bağlantıyı test eder
func (r *RedisCache) Ping() error {
	return r.client.Ping(r.ctx).Err()
}
