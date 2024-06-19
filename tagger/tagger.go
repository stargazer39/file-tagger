package tagger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/stargazer39/file-tagger/tagerror"
)

type Tagger struct {
	tagFile string
}

type TaggedFile struct {
	Name        string
	Tags        []string
	Description string
	Dir         bool
}

func NewTagger() Tagger {
	return Tagger{
		tagFile: ".tag",
	}
}

func (t *Tagger) SetCustomTagFile(tagFile string) {
	t.tagFile = tagFile
}

func (t *Tagger) ListFiles(rootDirPath string) (*[]TaggedFile, error) {
	db := NewMetadataDB(rootDirPath, t.tagFile)

	defer db.Close()

	taggedFile := []TaggedFile{}

	err2 := filepath.Walk(rootDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == rootDirPath || info.Name() == t.tagFile {
			return nil
		}

		if info.IsDir() || strings.HasPrefix(".", info.Name()) {
			return filepath.SkipDir
		}

		tags, err := db.GetTagsForFile(info.Name(), false)

		if err != nil {
			if !NoMetadata(err) {
				return err
			}
		}

		desc, err := db.GetDescriptionForFile(info.Name(), false)

		if err != nil {
			if !NoMetadata(err) {
				return err
			}
		}

		taggedFile = append(taggedFile, TaggedFile{
			Tags:        tags,
			Name:        info.Name(),
			Description: desc,
			Dir:         info.IsDir(),
		})
		return nil
	})

	if err2 != nil {
		return &taggedFile, err2
	}

	return &taggedFile, nil
}

func NoMetadata(err error) bool {
	if tErr, ok := err.(*tagerror.TagError); ok {
		if tErr.Is(tagerror.ErrNoMetadata) {
			return true
		}
	}

	return false
}

func (t *Tagger) TagFile(filePath string, tags []string) error {
	dir, file := filepath.Split(filePath)

	db := NewMetadataDB(dir, t.tagFile)

	defer db.Close()

	return db.TagFile(file, tags, true)
}

func (t *Tagger) SetDescriptionForFile(fileName string, desc string) error {
	dir, file := filepath.Split(fileName)

	db := NewMetadataDB(dir, t.tagFile)

	defer db.Close()

	return db.SetDescriptionForFile(file, desc, true)
}

func (t *Tagger) SetColor(fileName string, tags string) {

}
