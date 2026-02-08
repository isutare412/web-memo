package embedding

type Config struct {
	TEIBaseURL           string
	BM25BaseURL          string
	QdrantHost           string
	QdrantPort           int
	QdrantCollectionName string
	JobBufferSize        int
}
