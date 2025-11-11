package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Claims struct {
	UserID pgtype.UUID
	jwt.RegisteredClaims
}
