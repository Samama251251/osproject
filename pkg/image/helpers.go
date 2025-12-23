package image

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// LayerPath returns the filesystem path for a layer identified by digest.
func LayerPath(digest string) string {
	return filepath.Join(LyrDir, digest)
}

// RepoDir returns directory that contains the repository metadata file.
func RepoDir() string {
	return filepath.Dir(RepoFile)
}

// EnsureImageDirs creates required directories for images and layers if missing.
// It is safe to call multiple times (idempotent).
func EnsureImageDirs() error {
	if err := os.MkdirAll(LyrDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(RepoDir(), 0755); err != nil {
		return err
	}
	return nil
}

// AtomicWriteRepo writes repository metadata to `RepoFile` atomically by writing
// to a temp file in the same directory and renaming it into place.
func AtomicWriteRepo(data []byte, perm os.FileMode) error {
	dir := RepoDir()
	tmp, err := ioutil.TempFile(dir, ".firaaq-repo-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, RepoFile)
}

// ShortID returns first 12 characters of hash when available.
func ShortID(hash string) string {
	if len(hash) >= 12 {
		return hash[:12]
	}
	return hash
}
