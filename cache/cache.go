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

func (r *RedisCache) GetTweets2lark(id string) (isSend bool, err error) {
	key := fmt.Sprintf("tweetslark:%s", id)

	if ret, err := r.Red.Get(key).Result(); err != nil {
		return true, err
	} else {
		fmt.Println("1111redis get result: ", ret)
		return false, nil
	}
}
func (r *RedisCache) SetTweets2lark(id string) error {
	key := fmt.Sprintf("tweetslark:%s", id)
	ret := r.Red.Set(key, 1, 0)
	if err := ret.Err(); err != nil {
		return fmt.Errorf("redis set err: %s", err.Error())
	}
	return nil
}
