package archive

// Extractor is implemented by archive readers that can expand into dst directories.
type Extractor interface {
	Extract(dst string) error
}
