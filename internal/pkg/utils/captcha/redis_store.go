package captcha

import (
	"context"
	"github.com/mojocn/base64Captcha"
	"time"
)

type RedisStore struct {
	Expiration time.Duration
	PreKey     string
	Context    context.Context
}

func NewDefaultRedisStore() *RedisStore {
	return &RedisStore{
		Expiration: time.Second * 1800,
		PreKey:     "CAPTCHA_",
		Context:    context.TODO(),
	}
}

func (rs *RedisStore) UseWithCtx(ctx context.Context) base64Captcha.Store {
	rs.Context = ctx
	return rs
}

func (rs *RedisStore) Set(id string, value string) error {
	return nil
}

func (rs *RedisStore) Get(id string, clear bool) string {

	return ""
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {

	return true
}
