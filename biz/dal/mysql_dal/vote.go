package mysql_dal

import (
	"context"
	"fmt"
)

type UserTicket struct {
	ID       int64  `db:"id"`
	UserName string `db:"user_name"`
	Count    int64  `db:"count"`
}

func QueryVote(ctx context.Context, userName string) (int64, error) {
	query := "select id,user_name, count from user_ticket where user_name = ?"
	var u UserTicket
	err := conn.QueryRowCtx(ctx, &u, query, userName)
	if err != nil {
		return 0, err
	}
	return u.Count, nil
}

func AddVote(ctx context.Context, userNames []string) error {
	query := "insert into user_ticket (user_name, count) values %s on duplicate key update count = count + 1"
	temp := ""
	args := []any{}
	for i := 0; i < len(userNames); i++ {
		args = append(args, userNames[i])
		temp += "(?, 1)"
		if i != len(userNames)-1 {
			temp += ","
		}
	}
	query = fmt.Sprintf(query, temp)
	_, err := conn.ExecCtx(ctx, query, args...)
	return err

}
