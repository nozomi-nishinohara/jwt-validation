package infrastructure_test

import (
	"context"
	"testing"

	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
	"github.com/nozomi-nishinohara/jwt_validation/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestRedisCluster_1OK(t *testing.T) {
	redisCli := infrastructure.NewRedis()
	ctx := context.TODO()
	err := redisCli.Save(ctx, "aaa", &model.JSONWebKeys{
		Kid: "test",
	})
	jwk, err := redisCli.Get(ctx, "aaa")
	assert.NoError(t, err)
	if jwk != nil {
		assert.Equal(t, jwk.Kid, "test")
	}
}

func TestRedisCluster_2NG(t *testing.T) {
	redisCli := infrastructure.NewRedis()
	ctx := context.TODO()
	_, err := redisCli.Get(ctx, "test")
	assert.Equal(t, err, infrastructure.ErrNotFound)
}
