package ticket

import (
	"context"
	"flag"
	"graphql_vote/biz/dal/redis_dal"
	"time"

	"log"

	"github.com/libi/dcron"
	"github.com/libi/dcron/cron"
	"github.com/libi/dcron/driver"
	"github.com/redis/go-redis/v9"
)

const (
	DriverType_REDIS = "redis"
	DriverType_ETCD  = "etcd"
)

var (
	addr       = flag.String("addr", "127.0.0.1:6379", "the addr of driver service")
	serverName = flag.String("server_name", "server", "the server name of dcron in this process")
)

func InitTicket() {
	redisCli := redis.NewClient(&redis.Options{
		Addr: *addr,
	})
	drv := driver.NewRedisDriver(redisCli)
	dcron := dcron.NewDcron(
		*serverName,
		drv,
		cron.WithSeconds(),
	)
	// 2s执行一次
	err := dcron.AddFunc("create_ticket", "*/2 * * * * *", func() {
		err := redis_dal.SetTicket(context.Background())
		if err != nil {
			log.Printf("create ticket failed, %s", time.Now().Format("15:04:05.000"))
			return
		}
		log.Printf("create ticket success, %s", time.Now().Format("15:04:05.000"))
	})
	if err != nil {
		panic(err)
	}
	dcron.Start()
}
