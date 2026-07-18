package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/isutare412/web-memo/api/internal/config"
	"github.com/isutare412/web-memo/api/internal/core/enum"
)

var _ = Describe("Config", func() {
	Context("ToEmbeddingConfig", func() {
		It("maps the TEI model", func() {
			var cfg config.Config
			cfg.Embedding.TEIModel = enum.TEIModelQwen3Embedding4B

			Expect(cfg.ToEmbeddingConfig().TEIModel).To(Equal(enum.TEIModelQwen3Embedding4B))
		})
	})
})
