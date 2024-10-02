package aws

func chunkSlice[T any](s []T, size int) [][]T {
	var chunks [][]T
	for size < len(s) {
		s, chunks = s[size:], append(chunks, s[0:size:size])
	}
	chunks = append(chunks, s)
	return chunks
}
