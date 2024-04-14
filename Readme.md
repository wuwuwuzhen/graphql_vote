可以实现：

1. 分布式生成ticket



分布式生成ticket

1. 使用dcron库实现：使用 redis/etcd 同步服务节点列表及存活状态，在节点列表内使用一致性hash，选举可执行任务的节点。
2. 选举出的节点每间隔2s执行SetTicket任务
   ```
   // 在 Redis 中存储一个2 秒过期时间的ticket uuid
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

1. 通过各个节点在定时任务内抢锁方式实现，需要依赖各个节点系统时间完全一致，
