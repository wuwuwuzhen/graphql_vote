package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"graphql_vote/biz/dal/mysql_dal"
	"graphql_vote/biz/dal/redis_dal"
	"log"
	"time"
)

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, usernames []string, ticket string) (string, error) {
	lock := redis_dal.Lock
	timeout := time.After(1 * time.Second)
	tick := time.Tick(100 * time.Millisecond) 
	for {
		select {
		case <-timeout:
			log.Fatal("未能获取锁，已超时")
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
					return "ticket has expired", nil
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

// QueryVote is the resolver for the query_vote field.
func (r *queryResolver) QueryVote(ctx context.Context, username string) (int, error) {
	count, err := redis_dal.GetVote(ctx, username)
	if err == nil {
		return int(count), nil
	}
	count, err = mysql_dal.QueryVote(ctx, username)
	if err != nil {
		return 0, err
	}
	_ = redis_dal.SetVote(ctx, username, count)
	return int(count), nil
}

// GetTicket is the resolver for the get_ticket field.
func (r *queryResolver) GetTicket(ctx context.Context) (string, error) {
	ticket, err := redis_dal.GetTicket(ctx)
	if err != nil {
		return "", err
	}
	return ticket, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }