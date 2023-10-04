package redis

type Config struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}
