package enum

import "fmt"

type TEIModel string

const (
	TEIModelQwen3Embedding0_6B TEIModel = "qwen3-embedding-0.6b"
	TEIModelQwen3Embedding4B   TEIModel = "qwen3-embedding-4b"
)

func (m TEIModel) Validate() error {
	switch m {
	case TEIModelQwen3Embedding0_6B, TEIModelQwen3Embedding4B:
		return nil
	}
	return fmt.Errorf("unknown TEI model %q; valid values are %q or %q",
		string(m), TEIModelQwen3Embedding0_6B, TEIModelQwen3Embedding4B)
}
