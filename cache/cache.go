package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/scorpiotzh/mylog"
)

type RedisCache struct {
	Ctx context.Context
	Red *redis.Client
}

var (
	log = mylog.NewLogger("cache", mylog.LevelDebug)
)

func Initialize(red *redis.Client) *RedisCache {
	return &RedisCache{Red: red}
}

func (r *RedisCache) GetTweets2lark(id string) bool {
	key := fmt.Sprintf("tweetslark:%s", id)

	ret := r.Red.Get(key).Val()

	if ret == "1" {
		return true
	} else {
		return false
	}
}
func (r *RedisCache) SetTweets2lark(id string) error {
	key := fmt.Sprintf("tweetslark:%s", id)
	ret := r.Red.Set(key, 1, 0)
	if err := ret.Err(); err != nil {
		return fmt.Errorf("get coupon lock: redis set nx-->%s", err.Error())
	}
	return nil
}
