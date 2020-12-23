package infrastructure_test

import (
	"context"
	"testing"

	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
	"github.com/nozomi-nishinohara/jwt_validation/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestInMemory_1OK(t *testing.T) {
	inmemory := infrastructure.NewInMemory()
	c := context.TODO()
	inmemory.Save(c, "test", &model.JSONWebKeys{
		Kid: "test",
	})

	m, _ := inmemory.Get(c, "test")
	assert.Equal(t, m.Kid, "test")
}

func TestInMemory_2NG(t *testing.T) {
	inmemory := infrastructure.NewInMemory()
	c := context.TODO()
	m, err := inmemory.Get(c, "test")
	assert.Equal(t, err, infrastructure.ErrNotFound)
	assert.Nil(t, m)
}
