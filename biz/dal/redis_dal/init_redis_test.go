package redis_dal

import (
	"context"
	"fmt"
	"testing"

	"github.com/zeromicro/go-zero/core/logc"
)

func TestInitRedis(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRedis()
			ctx := context.Background()
			err := rds.SetCtx(ctx, "key", "hello world")
			if err != nil {
				logc.Error(ctx, err)
			}

			v, err := rds.GetCtx(ctx, "key")
			if err != nil {
				logc.Error(ctx, err)
			}
			fmt.Println(v)
		})
	}
}
