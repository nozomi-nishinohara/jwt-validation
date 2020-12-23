package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nozomi-nishinohara/jwt_validation/internal"
	"gopkg.in/yaml.v2"
)

var _setting *Setting

func GetSetting() *Setting {
	if _setting == nil {
		setting := loadSetting()
		_setting = setting
		if _setting.Cache.Name == "" {
			_setting.Cache.Name = InMemory
			_setting.Cache.Time = 30
		}
		if len(_setting.Oauths) == 0 {
			panic(errors.New("Oauth Providerが設定されていません。"))
		}
	}
	return _setting
}

func loadSetting() *Setting {
	setting := &Setting{}
	file := internal.Getenv("VALIDATION_FILE_NAME", "oauth")
	if _, err := os.Stat(file + ".yaml"); !os.IsNotExist(err) {
		filename := file + ".yaml"
		if buf, err := ioutil.ReadFile(filename); err != nil {
			panic(err)
		} else {
			expandContent := []byte(os.ExpandEnv(string(buf)))
			if err := yaml.Unmarshal(expandContent, setting); err != nil {
				panic(err)
			}
			return setting
		}
	}

	if _, err := os.Stat(file + ".json"); !os.IsNotExist(err) {
		filename := file + ".json"
		if buf, err := ioutil.ReadFile(filename); err != nil {
			panic(err)
		} else {
			expandContent := []byte(os.ExpandEnv(string(buf)))
			if err := json.Unmarshal(expandContent, setting); err != nil {
				panic(err)
			}
			return setting
		}
	}

	return nil
}

func (s *Setting) GetOauth(iss string) *Oauth {
	for _, oauth := range s.Oauths {
		if oauth.Iss == iss {
			return oauth
		}
	}
	return nil
}

type (
	Oauth struct {
		Domain    string   `yaml:"domain" json:"domain"`
		Iss       string   `yaml:"iss" json:"iss"`
		Aud       []string `yaml:"aud" json:"aud"`
		JwkSetUri string   `yaml:"jwk-set-uri" json:"jwk-set-uri"`
	}
	Setting struct {
		Oauths []*Oauth `yaml:"oauth" json:"oauth"`
		Cache  Cache    `yaml:"cache" json:"cache"`
	}
	Cache struct {
		Name CacheName `yaml:"name" json:"name"`
		Time int       `yaml:"time" json:"time"`
	}
)

func (c *Cache) GetTime() int64 {
	if c.Time == 0 {
		c.Time = 30 // default expiration of 30 minute
	}
	return time.Now().Add(time.Duration(c.Time) * time.Minute).UnixNano()
}

type CacheName string

const (
	InMemory = CacheName("inmemory")
	Redis    = CacheName("redis")
)

func (c CacheName) String() string {
	switch c {
	case InMemory:
		return "inmemory"
	case Redis:
		return "redis"
	default:
		return string(c)
	}
}

func (c CacheName) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CacheName) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var cache CacheName
	switch s {
	case "inmemmory":
		cache = InMemory
	case "redis":
		cache = Redis
	default:
		cache = CacheName(s)
	}
	*c = cache
	return nil
}
