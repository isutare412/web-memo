package embedding

import "github.com/isutare412/web-memo/api/internal/core/enum"

// Chunk limits are UTF-8 bytes, not runes. Each profile is sized so a full
// teiBatchSize batch stays within that deployment's TEI --max-batch-tokens
// budget without truncation (Korean text runs ~1 token per 3 bytes, English
// ~1 per 4).
type modelProfile struct {
	vectorSize        uint64
	maxChunkBytes     int
	chunkOverlapBytes int
}

var modelProfiles = map[enum.TEIModel]modelProfile{
	enum.TEIModelQwen3Embedding0_6B: {vectorSize: 1024, maxChunkBytes: 4096, chunkOverlapBytes: 256},
	enum.TEIModelQwen3Embedding4B:   {vectorSize: 2560, maxChunkBytes: 8192, chunkOverlapBytes: 512},
}
