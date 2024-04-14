package redis_dal

import (
	"context"
	"fmt"
	"strconv"
)

func SetVote(ctx context.Context, userName string, count int64) error {
	err := rds.SetexCtx(ctx, userName, fmt.Sprintf("%d", count), 10)
	if err != nil {
		return err
	}
	return nil
}

func DeleteVote(ctx context.Context, userName string) error {
	_, err := rds.DelCtx(ctx, userName)
	if err != nil {
		return err
	}
	return nil
}

func GetVote(ctx context.Context, userName string) (int64, error) {
	exist, err := rds.ExistsCtx(ctx, userName)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, fmt.Errorf("key not exist")
	}
	v, err := rds.GetCtx(ctx, userName)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
