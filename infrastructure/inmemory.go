package infrastructure

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
	"github.com/nozomi-nishinohara/jwt_validation/domain/repository"
)

type (
	item struct {
		value   string
		expires int64
	}
	Cache struct {
		items map[string]*item
		mu    sync.Mutex
	}

	inmemory struct {
		data map[string]*model.JSONWebKeys
	}
)

func NewInMemory() repository.Cache {
	c := &Cache{items: make(map[string]*item)}
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				c.mu.Lock()
				for k, v := range c.items {
					if v.Expired(time.Now().UnixNano()) {
						log.Printf("%v has expires at %d", c.items, time.Now().UnixNano())
						delete(c.items, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}()
	return c
}

func (i *item) Expired(time int64) bool {
	if i.expires == 0 {
		return true
	}
	return time > i.expires
}

func (c *Cache) Get(ctx context.Context, key string) (*model.JSONWebKeys, error) {
	c.mu.Lock()
	var s string
	if v, ok := c.items[key]; ok {
		s = v.value
	}
	c.mu.Unlock()
	if s != "" {
		return model.NewJsonToJWKS([]byte(s))
	} else {
		return nil, ErrNotFound
	}
}

func (c *Cache) Save(ctx context.Context, key string, jwks *model.JSONWebKeys) error {
	c.mu.Lock()
	if _, ok := c.items[key]; !ok {
		c.items[key] = &item{
			value:   jwks.ToJson(),
			expires: model.GetSetting().Cache.GetTime(),
		}
	}
	c.mu.Unlock()
	return nil
}
