package embedding

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/qdrant/go-client/qdrant"

	"github.com/isutare412/web-memo/api/internal/core/enum"
)

var _ = Describe("Client", func() {
	Context("NewClient", func() {
		It("resolves the 0.6b profile", func() {
			client, err := NewClient(Config{
				TEIModel:   enum.TEIModelQwen3Embedding0_6B,
				QdrantHost: "localhost",
				QdrantPort: 6334,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(client.profile.vectorSize).To(Equal(uint64(1024)))
			Expect(client.profile.maxChunkBytes).To(Equal(4096))
			Expect(client.profile.chunkOverlapBytes).To(Equal(256))
		})

		It("resolves the 4b profile", func() {
			client, err := NewClient(Config{
				TEIModel:   enum.TEIModelQwen3Embedding4B,
				QdrantHost: "localhost",
				QdrantPort: 6334,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(client.profile.vectorSize).To(Equal(uint64(2560)))
			Expect(client.profile.maxChunkBytes).To(Equal(8192))
			Expect(client.profile.chunkOverlapBytes).To(Equal(512))
		})

		It("rejects an unknown model", func() {
			_, err := NewClient(Config{
				TEIModel:   enum.TEIModel("bogus"),
				QdrantHost: "localhost",
				QdrantPort: 6334,
			})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("validateDenseDimension", func() {
		buildInfo := func(size uint64) *qdrant.CollectionInfo {
			return &qdrant.CollectionInfo{
				Config: &qdrant.CollectionConfig{
					Params: &qdrant.CollectionParams{
						VectorsConfig: qdrant.NewVectorsConfigMap(map[string]*qdrant.VectorParams{
							"dense": {Size: size, Distance: qdrant.Distance_Cosine},
						}),
					},
				},
			}
		}

		It("accepts a matching dimension", func() {
			err := validateDenseDimension(buildInfo(2560), "memo_chunks", enum.TEIModelQwen3Embedding4B, 2560)
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects a mismatched dimension", func() {
			err := validateDenseDimension(buildInfo(1024), "memo_chunks", enum.TEIModelQwen3Embedding4B, 2560)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("1024"))
			Expect(err.Error()).To(ContainSubstring("2560"))
			Expect(err.Error()).To(ContainSubstring("memo_chunks"))
			Expect(err.Error()).To(ContainSubstring("set tei-model"))
		})

		It("rejects a collection without dense vector config", func() {
			info := &qdrant.CollectionInfo{}
			err := validateDenseDimension(info, "memo_chunks", enum.TEIModelQwen3Embedding4B, 2560)
			Expect(err).To(HaveOccurred())
		})
	})
})
