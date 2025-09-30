package config

type (
	AppSettings struct {
		PORT int
	}
	DbSettings struct {
		AutoConnect bool
		Host        string
		Port        string
		Database    string
		Username    string
		Password    string
		Dialect     string
	}
	RedisSettings struct{}
	OrionConfig   struct {
		App   AppSettings
		Db    DbSettings
		Redis RedisSettings
	}
)
