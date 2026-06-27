package document_driver_manager

import "time"

type Config struct {
	Driver      string
	DSN         string
	Database    string
	PingTimeout time.Duration
}
