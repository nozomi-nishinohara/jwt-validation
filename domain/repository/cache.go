package repository

import (
	"context"

	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
)

type Cache interface {
	Get(c context.Context, key string) (*model.JSONWebKeys, error)
	Save(c context.Context, key string, jwks *model.JSONWebKeys) error
}
