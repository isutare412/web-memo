package redis

type Config struct {
	Addr     string `koanf:"addr" validate:"required"`
	Password string `koanf:"password" validate:"required"`
}
