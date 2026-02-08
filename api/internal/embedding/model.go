package embedding

type embedRequest struct {
	Inputs []string `json:"inputs"`
}

type sparseVector struct {
	Indices []uint32
	Values  []float32
}

type sparseEmbedResponse struct {
	Indices []uint32  `json:"indices"`
	Values  []float32 `json:"values"`
}
