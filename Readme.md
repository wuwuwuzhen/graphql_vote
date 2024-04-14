## 分布式生成随机ticket

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

1. 通过各个节点在定时任务内抢锁方式实现，需要依赖各个节点系统时间完全一致。

## 投票

1. 使用redis分布式锁保证数据一致性，在获取分布式锁之前会不断重试。
2. 在锁被成功获取后，通过 `redis_dal.GetTicket(ctx)` 从 Redis 中获取当前有效的票据。如果票据不匹配或已经失效（例如，过期或被更新）返回"ticket has expired".
