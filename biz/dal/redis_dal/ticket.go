package redis_dal

import (
	"context"

	"github.com/google/uuid"
)

var ticketKey = "graphql_ticket"

func SetTicket(ctx context.Context) error {
	newUUID := uuid.New()
	err := rds.SetexCtx(ctx, ticketKey, newUUID.String(), 2)
	if err != nil {
		return nil
	}
	return nil
}

func GetTicket(ctx context.Context) (string, error) {
	v, err := rds.GetCtx(ctx, ticketKey)
	if err != nil {
		return "", err
	}
	return v, nil
}
