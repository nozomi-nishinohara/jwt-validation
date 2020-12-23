package infrastructure

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
	"github.com/nozomi-nishinohara/jwt_validation/domain/repository"
	"github.com/nozomi-nishinohara/jwt_validation/internal"
)

type (
	client interface {
		Get(ctx context.Context, key string) *redis.StringCmd
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	}
	redisClient struct {
		cmd client
	}
)

func NewRedis() repository.Cache {
	var client client
	endpoint := internal.Getenv("REDIS_CLUSTER_ENDPOINT", "redis") + ":" + internal.Getenv("REDIS_CLUSTER_PORT", "6379")
	ctx := context.Background()
	RedisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{endpoint},
	})
	err := RedisClusterClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		RedisClusterClient = nil
		RedisClient := redis.NewClient(&redis.Options{
			Addr:     endpoint,
			Password: "",
			DB:       0,
		})
		if err = RedisClient.Ping(ctx).Err(); err != nil {
			panic(err)
		}
		client = RedisClient
	} else {
		client = RedisClusterClient
	}
	if client == nil {
		panic("No Support Redis")
	}
	return &redisClient{client}
}

func (r *redisClient) Get(c context.Context, key string) (*model.JSONWebKeys, error) {
	if res, err := r.cmd.Get(c, key).Result(); err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	} else {
		return model.NewJsonToJWKS([]byte(res))
	}
}

func (r *redisClient) Save(c context.Context, key string, jwks *model.JSONWebKeys) error {
	expiration := model.GetSetting().Cache.GetTime()
	var err error
	if _, err = r.cmd.Set(c, key, []byte(jwks.ToJson()), time.Duration(expiration)).Result(); err != nil {
		return err
	}
	return nil
}
