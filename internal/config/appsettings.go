package config

type (
	AppSettings struct {
		PORT int
	}
	DbSettings    struct{}
	RedisSettings struct{}
	OrionConfig   struct {
		App   AppSettings
		Db    DbSettings
		Redis RedisSettings
	}
)
