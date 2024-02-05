package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host             string
	Port             uint16
	ConnectTimeout   time.Duration
	QueryTimeout     time.Duration
	Username         string
	Password         string
	DBName           string
	MigrationVersion int64
}

func (c Config) connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.DBName)
}
