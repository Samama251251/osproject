package container

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// IsHex reports whether s contains only hexadecimal characters.
func IsHex(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			continue
		case c >= 'a' && c <= 'f':
			continue
		case c >= 'A' && c <= 'F':
			continue
		default:
			return false
		}
	}
	return true
}

// IsValidDigestPrefix validates that prefix can be used to identify a digest
// present on disk (hex characters, length between 1 and DigestStdLen).
func IsValidDigestPrefix(prefix string) bool {
	if prefix == "" || len(prefix) > DigestStdLen {
		return false
	}
	return IsHex(prefix)
}

// ContainerPath returns the absolute path for a container with given digest.
func ContainerPath(digest string) string {
	return filepath.Join(containerPath, digest)
}

// ContainerMountPath returns the mount directory path for the container.
func ContainerMountPath(digest string) string {
	return filepath.Join(ContainerPath(digest), "mnt")
}

// EnsureContainerDirs creates the common container directories (container dir
// and network namespace dir) with safe permissions. It is idempotent.
func EnsureContainerDirs(digest string) error {
	if strings.TrimSpace(digest) == "" {
		return os.ErrInvalid
	}
	if err := os.MkdirAll(ContainerMountPath(digest), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(containerNetNsPath, digest), 0700); err != nil {
		return err
	}
	return nil
}

// AtomicWriteFile writes data to a temporary file in the same directory and
// renames it into place. This reduces the window where a partial file may be
// observed by other readers.
func AtomicWriteFile(filename string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	tmp, err := ioutil.TempFile(dir, ".firaaq-write-*")
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
	return os.Rename(tmpName, filename)
}

// ShortID returns the 12-character short form of digest if available.
func ShortID(digest string) string {
	if len(digest) >= 12 {
		return digest[:12]
	}
	return digest
}
