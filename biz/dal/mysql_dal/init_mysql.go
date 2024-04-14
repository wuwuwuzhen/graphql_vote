package mysql_dal

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var conn sqlx.SqlConn

func InitMysql() {
	dsn := "april:1001@tcp(127.0.0.1:3306)/graphql_ticket?charset=utf8mb4&parseTime=True&loc=Local"
	conn = sqlx.NewMysql(dsn)
}
