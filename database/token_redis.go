package database

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	rfsKey := redisKey(arg.UserId)

	trans := client.TxPipeline()
	trans.HSet(r.Context(), key, arg)
	trans.ExpireAt(r.Context(), key, arg.Expiry.UTC())
	trans.Set(r.Context(), rfsKey, arg.RefreshToken, 0)
	trans.ExpireAt(r.Context(), rfsKey, arg.Expiry.UTC())

	_, err := trans.Exec(r.Context())
	return err
}

func GetToken(r *http.Request, client *redis.Client, refreshToken string) (RedisToken, error) {
	key := redisKey(refreshToken)
	v := RedisToken{}
	var getAll *redis.MapStringStringCmd
	var expTime *redis.DurationCmd

	_, err := client.TxPipelined(r.Context(), func(p redis.Pipeliner) error {
		getAll = client.HGetAll(r.Context(), key)
		expTime = client.ExpireTime(r.Context(), key)
		return nil
	})
	if err != nil {
		return v, err
	}

	getAll.Scan(&v)
	v.Expiry = time.Unix(int64(expTime.Val().Seconds()), 0).UTC()
	v.RefreshToken = refreshToken
	return v, err
}

func DeleteToken(r *http.Request, client *redis.Client, arg pgtype.UUID) error {
	rfsKey := redisKey(arg.String())
	rfsToken, err := client.Get(r.Context(), rfsKey).Result()
	if err != nil {
		return err
	}
	key := redisKey(rfsToken)

	err = client.Del(r.Context(), rfsKey, key).Err()
	return err
}
