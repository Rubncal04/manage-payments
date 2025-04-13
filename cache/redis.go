package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, port string, db int) (*RedisCache, error) {
	fmt.Println(addr + ":" + port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: password,
		DB:       db,
	})

	// Verificar conexión
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

// Set almacena un valor en caché
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, json, expiration).Err()
}

// Get obtiene un valor de caché
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete elimina un valor de caché
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Clear limpia toda la caché
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

// GenerateKey genera una clave única para caché
func GenerateKey(prefix string, params ...any) string {
	key := prefix
	for _, param := range params {
		key += ":" + toString(param)
	}
	return key
}

// toString convierte cualquier valor a string
func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
