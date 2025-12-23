package image

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

const (
	RepoFile = "/var/lib/firaaq/images/repositories.json"
	LyrDir   = "/var/lib/firaaq/images/layers"
)

type Image struct {
	v1.Image
	ID         string
	Registry   string
	Repository string
	Name       string
	Tag        string
}

// NewImage returns a new Image.
//
// It pulls image ans sets the Registry, Repository, Name, and Tag.
func NewImage(src string) (*Image, error) {
	tag, err := name.NewTag(src)
	if err != nil {
		return nil, err
	}
	// Specify platform to match host architecture (arm64 for M4 Mac)
	platform := v1.Platform{
		Architecture: runtime.GOARCH,
		OS:           runtime.GOOS,
	}
	img, err := crane.Pull(tag.Name(), crane.WithPlatform(&platform))
	if err != nil {
		return nil, err
	}
	digest, err := img.Digest()
	if err != nil {
		return nil, err
	}
	newImage := &Image{
		Image:      img,
		ID:         digest.Hex,
		Registry:   tag.RegistryStr(),
		Repository: tag.RepositoryStr(),
		Name:       tag.Name(),
		Tag:        tag.TagStr(),
	}
	return newImage, nil
}

// Exists checks for image existence in local storage.
func (i *Image) Exists() (bool, error) {
	// Build current repository cache before checking for this digest.
	images, err := GetAll()
	if err != nil {
		return false, err
	}
	for _, img := range images {
		if img.ID == i.ID {
			return true, nil
		}
	}
	return false, nil
}

func GetAll() ([]*Image, error) {
	repos := make(Repositories)
	imgs := make([]*Image, 0)

	data, err := ioutil.ReadFile(RepoFile)
	if os.IsNotExist(err) {
		return imgs, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}

	// Expand repositories.json entries into Image structs so callers can inspect them.
	for repo, image := range repos {
		for nameTag, hash := range image {
			newImg := &Image{
				ID:         strings.TrimLeft(hash, "sha256:"),
				Repository: repo,
				Tag:        strings.Split(nameTag, ":")[1],
			}
			imgs = append(imgs, newImg)
		}
	}
	return imgs, nil
}
