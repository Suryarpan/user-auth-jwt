package database

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RefreshNameSpace string = "user_auth:user_refresh"
)

func redisKey(k string) string {
	return RefreshNameSpace + ":" + k
}

func StoreToken(r *http.Request, client *redis.Client, arg RedisToken) error {
	key := redisKey(arg.RefreshToken)
	_, err := client.HSet(r.Context(), key, arg).Result()
	if err != nil {
		return err
	}
	_, err = client.ExpireAt(r.Context(), key, arg.Expiry.UTC()).Result()
	return err
}

func GetToken(r *http.Request, client *redis.Client, refreshToken string) (RedisToken, error) {
	key := redisKey(refreshToken)
	v := RedisToken{}
	err := client.HGetAll(r.Context(), key).Scan(&v)
	if err != nil {
		return v, err
	}
	d, err := client.ExpireTime(r.Context(), key).Result()
	if err != nil {
		return v, err
	}
	v.Expiry = time.Unix(int64(d.Seconds()), 0).UTC()
	v.RefreshToken = refreshToken
	return v, err
}
