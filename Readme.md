# Tiny Vote

运行 `go run server.go`

## 遵循GraphQl的接口设计

```type
type Mutation {
    vote(usernames: [String!]!, ticket: String!): String!
}

type Query {
    query(username: String!): Int!
    cas: String!
}

```

## 分布式生成随机ticket

1. 使用dcron库实现：使用 redis/etcd 同步服务节点列表及存活状态，在节点列表内使用一致性hash，选举可执行任务的节点。
2. 每间隔2s选举节点执行SetTicket任务

   ```
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
   ```
3. 在 Redis 中存储一个2 秒过期时间的ticket uuid

   ```
   var ticketKey = "graphql_ticket"
   func SetTicket(ctx context.Context) error {
   	newUUID := uuid.New()
   	err := rds.SetexCtx(ctx, ticketKey, newUUID.String(), 2)
   	if err != nil {
   		return nil
   	}
   	return nil
   }
   ```

为什么不用分布式锁实现生成随机ticket：

1. 通过各个节点在定时任务内抢锁方式实现，需要依赖各个节点系统时间完全一致。

## 投票

1. 使用redis分布式锁保证数据一致性，轮询1s，超时则返回"未能获取锁，已超时"。
2. 在锁被成功获取后，通过 `redis_dal.GetTicket(ctx)` 从 Redis 中获取当前有效的票据。如果ticket失效，返回"ticket 失效"。若ticket合法，将指定的（多个）⽤户的票数+1。
3. 投票数据用mysql持久化，写入mysql后删除对应缓存。

   ```
   func (r *mutationResolver) Vote(ctx context.Context, usernames []string, ticket string) (string, error) {
   	lock := redis_dal.Lock
   	timeout := time.After(1 * time.Second)
   	tick := time.Tick(100 * time.Millisecond) 
   	for {
   		select {
   		case <-timeout:
   			log.Print("未能获取锁，已超时")
   			return "未能获取锁，已超时", nil
   		case <-tick:
   			acquire, err := lock.AcquireCtx(ctx)
   			switch {
   			case err != nil:
   				log.Fatal("尝试获取锁时发生错误", err)
   				return "尝试获取锁时发生错误", err
   			case acquire:
   				log.Println("获取到锁")
   				defer lock.Release() 
   				curTicket, err := redis_dal.GetTicket(ctx)
   				if err != nil {
   					return "", err
   				}
   				if ticket != curTicket {
   					return "ticket 失效", nil
   				}

   				err = mysql_dal.AddVote(ctx, usernames)
   				if err != nil {
   					return "", err
   				}
   				for _, username := range usernames {
   					_ = redis_dal.DeleteVote(ctx, username)
   				}
   				return "成功", nil
   			case !acquire:
   				log.Println("未获取到锁，等待中")
   			}
   		}
   	}
   }
   ```

## 查询

1. 以用户名为key读缓存，如果缓存中key不存在，则访问数据库后将数据载入缓存。
2. singlefight防止缓存击穿（暂未实现，来不及了QAQ）
