package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStorage implementa LimiterStorage usando Redis.
type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorage cria uma nova instância de RedisStorage.
func NewRedisStorage(addr string) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Incrementa o contador de requisições para uma chave e define o tempo de expiração da janela.
func (r *RedisStorage) Increment(key string, window time.Duration) (int, error) {
	countKey := "count:" + key
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(r.ctx, countKey)
	pipe.Expire(r.ctx, countKey, window)
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return 0, err
	}
	return int(incr.Val()), nil
}

// Retorna o número atual de requisições para a chave.
func (r *RedisStorage) GetCount(key string) (int, error) {
	countKey := "count:" + key
	val, err := r.client.Get(r.ctx, countKey).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Bloqueia a chave por um determinado tempo.
func (r *RedisStorage) Block(key string, duration time.Duration) error {
	blockKey := "block:" + key
	return r.client.Set(r.ctx, blockKey, time.Now().Add(duration).Unix(), duration).Err()
}

// Verifica se a chave está bloqueada e retorna até quando.
func (r *RedisStorage) IsBlocked(key string) (bool, time.Time, error) {
	blockKey := "block:" + key
	val, err := r.client.Get(r.ctx, blockKey).Result()
	if err == redis.Nil {
		return false, time.Time{}, nil
	}
	if err != nil {
		return false, time.Time{}, err
	}
	ts, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, time.Time{}, err
	}
	until := time.Unix(ts, 0)
	if time.Now().Before(until) {
		return true, until, nil
	}
	return false, time.Time{}, nil
}

// Reseta o contador e o bloqueio da chave.
func (r *RedisStorage) Reset(key string) error {
	countKey := "count:" + key
	blockKey := "block:" + key
	return r.client.Del(r.ctx, countKey, blockKey).Err()
}
