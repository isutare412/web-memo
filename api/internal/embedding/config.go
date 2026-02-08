package embedding

type Config struct {
	TEIBaseURL              string
	QdrantHost              string
	QdrantPort              int
	QdrantCollectionName    string
	JobBufferSize           int
	MinSearchScoreThreshold float32
	MaxSearchResults        int
}
