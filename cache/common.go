package cache

import (
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// RedisClient 包含一个 *redis.Client 成员和一个互斥锁
type redisClient struct {
	Client *redis.Client
}

var RedisClient *redisClient

// InitCache 在中间件中初始化 Redis 连接，防止循环导包，所以放在这里
func InitCache() {
	Redis()
}

// Redis 在中间件中初始化 Redis 连接
func Redis() {
	db, _ := strconv.ParseUint("2", 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       int(db),
	})
	_, err := client.Ping().Result() // 心跳检测
	if err != nil {
		logging.Info(err)
		panic(err)
	}

	// 初始化 RedisClient
	RedisClient = &redisClient{
		Client: client,
	}
}
func SetHashFieldWithExpiration(key string, field string, value interface{}, expiration time.Duration) error {
	// 设置哈希表中的字段
	if err := RedisClient.Client.HSet(key, field, value).Err(); err != nil {
		return err
	}
	// 为整个哈希表设置过期时间
	if _, err := RedisClient.Client.Expire(key, expiration).Result(); err != nil {
		return err
	}
	return nil
}
