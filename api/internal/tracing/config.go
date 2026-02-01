package tracing

type Config struct {
	Enabled          bool    `koanf:"enabled"`
	ServiceName      string  `koanf:"service-name" validate:"required"`
	SamplingRatio    float64 `koanf:"sampling-ratio" validate:"gte=0,lte=1"` // 0.0 - 1.0
	OTLPGRPCEndpoint string  `koanf:"otlp-grpc-endpoint" validate:"required"`
}
