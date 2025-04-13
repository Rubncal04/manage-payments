package cache

import (
	"context"
	"time"
)

// Cache define la interfaz que cualquier implementación de caché debe seguir
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// CacheableRepository define la interfaz que cualquier repositorio que use caché debe implementar
type CacheableRepository interface {
	GetCache() Cache
	SetCache(cache Cache)
}

// BaseCacheableRepository proporciona una implementación base para repositorios con caché
type BaseCacheableRepository struct {
	cache Cache
}

func (r *BaseCacheableRepository) GetCache() Cache {
	return r.cache
}

func (r *BaseCacheableRepository) SetCache(cache Cache) {
	r.cache = cache
}

// WithCache es un helper para manejar operaciones con caché
func WithCache[T any](ctx context.Context, cache Cache, key string, dest T, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	if cache != nil {
		// Intentar obtener de caché
		err := cache.Get(ctx, key, &dest)
		if err == nil {
			return dest, nil
		}
	}

	// Si no está en caché o hay error, obtener de la fuente original
	result, err := fetchFunc()
	if err != nil {
		return dest, err
	}

	// Guardar en caché si está disponible
	if cache != nil {
		if err := cache.Set(ctx, key, result, ttl); err != nil {
			// Log error pero no fallar la operación
			// TODO: Agregar logger apropiado
		}
	}

	return result, nil
}

// InvalidateCache es un helper para invalidar caché después de operaciones de escritura
func InvalidateCache(ctx context.Context, cache Cache, keys ...string) {
	if cache == nil {
		return
	}

	for _, key := range keys {
		if err := cache.Delete(ctx, key); err != nil {
			// Log error pero no fallar la operación
			// TODO: Agregar logger apropiado
		}
	}
}
