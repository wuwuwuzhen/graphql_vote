package redis_dal

import "github.com/zeromicro/go-zero/core/stores/redis"

var (
	rds  *redis.Redis
	Lock *redis.RedisLock
)

func InitRedis() {
	var err error
	conf := redis.RedisConf{
		Host: "127.0.0.1:6379",
		Type: "node",
		Tls:  false,
	}

	rds, err = redis.NewRedis(conf)
	if err != nil {
		panic(err)
	}
	// 初始化分布式锁
	Lock = redis.NewRedisLock(rds, "vote")
	// 设置过期时间
	Lock.SetExpire(1)
}
