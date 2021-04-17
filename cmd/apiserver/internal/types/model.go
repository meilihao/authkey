package types

import (
	"time"
)

const (
	RoleNone = iota
	RoleAdmin
)

type ClientGroup struct {
	Id        int64 `layer:";pk;autoincr"`
	Name      string
	CreatedAt time.Time `layer:";created_at"`
	UpdatedAt time.Time `layer:";updated_at"`
}
