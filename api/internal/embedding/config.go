package embedding

import "github.com/isutare412/web-memo/api/internal/core/enum"

type Config struct {
	TEIModel             enum.TEIModel
	TEIBaseURL           string
	BM25BaseURL          string
	QdrantHost           string
	QdrantPort           int
	QdrantCollectionName string
	JobBufferSize        int
}
