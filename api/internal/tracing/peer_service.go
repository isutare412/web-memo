package tracing

import semconv "go.opentelemetry.io/otel/semconv/v1.40.0"

var (
	PeerServiceGoogleOIDC = semconv.ServicePeerName("google-oidc")
	PeerServicePostgres   = semconv.ServicePeerName("postgres")
	PeerServiceRedis      = semconv.ServicePeerName("redis")
	PeerServiceInternet   = semconv.ServicePeerName("internet")
	PeerServiceTEI        = semconv.ServicePeerName("huggingface-tei-qwen3-embedding")
	PeerServiceBM25       = semconv.ServicePeerName("huggingface-bm25")
	PeerServiceQdrant     = semconv.ServicePeerName("qdrant")
)
