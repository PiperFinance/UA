package models

import (
	"database/sql"
	"github.com/mileusna/useragent"
	"golang.org/x/net/context"
)

type Country string
type IP sql.NullString

func (ip IP) Country(ctx context.Context) Country {
	panic("Has not Yet Been implemented :) ")
	return ""
}

// UserAgent is device definition in request header
type UserAgent string

func (agent UserAgent) Parse() useragent.UserAgent {
	return useragent.Parse(string(agent))
}
