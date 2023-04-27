package dbmodels

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID       string `bun:"id,pk,autoincrement"`
	Login    string `bun:"login,notnull"`
	Password string `bun:"password,notnull"`
	Role     string `bun:"role,notnull"`
}
