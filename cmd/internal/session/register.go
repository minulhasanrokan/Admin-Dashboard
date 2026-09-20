package session

import (
	"encoding/gob"
	"time"

	"admin-dashboard/cmd/internal/dtos"
)

func init() {
	gob.Register(dtos.RightInfo{})
	gob.Register(time.Time{})
}
