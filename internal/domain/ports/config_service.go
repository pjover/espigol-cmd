package ports

import "time"

type ConfigService interface {
	GetString(key string) string
	SetString(key string, value string) error
	GetTime(key string) time.Time
	SetTime(key string, value time.Time) error
	Init()
}
