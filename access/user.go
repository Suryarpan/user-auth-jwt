package access

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
)

type (
	UserSearch struct {
		PvtId  int
		UserId pgtype.UUID
		UserNm string
	}

	UserOper interface {
		CreateUser(ctx context.Context, lo *slog.Logger, arg UserPublic) (User, error)
		GetUser(ctx context.Context, lo *slog.Logger, arg UserSearch) (User, error)
		UpdateUser(ctx context.Context, lo *slog.Logger, arg UserPublic) (User, error)
		DeleteUser(ctx context.Context, lo *slog.Logger, arg UserSearch) (int, error)
	}

	TokenOper interface {
		StoreToken(ctx context.Context, lo *slog.Logger, arg RedisToken) error
		GetToken(ctx context.Context, lo *slog.Logger, arg string) (RedisToken, error)
		DeleteToken(ctx context.Context, lo *slog.Logger, arg pgtype.UUID) error
	}
)
