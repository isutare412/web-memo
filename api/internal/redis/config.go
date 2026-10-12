package redis

type Config struct {
	Addr      string `koanf:"addr" validate:"required"`
	Username  string `koanf:"username"`
	Password  string `koanf:"password" validate:"required"`
	KeyPrefix string `koanf:"key-prefix"`
}
