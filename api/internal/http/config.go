package http

type Config struct {
	Port int `mapstructure:"port" validate:"required"`
}
