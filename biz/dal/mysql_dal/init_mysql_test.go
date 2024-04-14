package mysql_dal

import (
	"context"
	"fmt"
	"testing"
)

func TestInitMysql(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitMysql()
			// resp, _ := QueryVote(context.Background(), "april")
			// fmt.Println(resp)
			err := AddVote(context.Background(), []string{"april", "test"})
			fmt.Println(err)
		})
		// r, err := conn.ExecCtx(context.Background(), "insert into user_ticket (user_name, count) values (?, ?)", "april", 100)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(r.RowsAffected())

	}
}
