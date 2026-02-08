package tracing

import semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

var (
	PeerServiceGoogleOIDC = semconv.PeerService("google-oidc")
	PeerServicePostgres   = semconv.PeerService("postgres")
	PeerServiceRedis      = semconv.PeerService("redis")
	PeerServiceInternet   = semconv.PeerService("internet")
	PeerServiceTEI        = semconv.PeerService("huggingface-tei-qwen3-embedding")
	PeerServiceQdrant     = semconv.PeerService("qdrant")
)
