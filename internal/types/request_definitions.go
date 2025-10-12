package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Claims struct {
	UserID pgtype.UUID
	jwt.RegisteredClaims
}

type ChangePasswordParams struct {
	OldPasswd string `json:"old_passwd"`
	NewPasswd string `json:"new_passwd"`
}
