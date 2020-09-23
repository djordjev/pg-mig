package filesystem

import (
	"errors"
	"github.com/spf13/afero"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var InvalidFileFormat = errors.New("invalid file format")

// CreateMigrationFile - creates a new file in path directory
func (fs *Filesystem) CreateMigrationFile(name string, location string) error {
	filename := filepath.Join(location, name)

	_, err := fs.Fs.Create(filename)
	if err != nil {
		return err
	}

	return nil
}

func (fs *Filesystem) GetFiles(from time.Time, to time.Time, direction Direction) ([]string, error) {
	if direction == DirectionUp {
		return fs.getUpFiles(from, to)
	}

	return nil, errors.New("invalid direction")
}

func (fs *Filesystem) getUpFiles(from time.Time, to time.Time) ([]string, error) {
	config, err := fs.LoadConfig()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, 10)

	files, err := afero.ReadDir(fs.Fs, config.Path)

	namePattern := regexp.MustCompile("^mig_([0-9]+)(.*)_up.sql$")
	tsPattern := regexp.MustCompile("^mig_([0-9]+).*_up.sql$")

	for i := 0; i < len(files); i++ {
		current := files[i]
		currentName := current.Name()

		matched := namePattern.FindString(currentName)
		if err != nil {
			return nil, err
		}

		if matched == "" {
			continue
		}

		ts := tsPattern.FindStringSubmatch(currentName)
		if ts == nil || len(ts) < 2 {
			return nil, InvalidFileFormat
		}

		strTs := ts[1]
		intTs, err := strconv.ParseInt(strTs, 10, 64)
		if err != nil {
			return nil, err
		}

		if intTs < from.Unix() || intTs > to.Unix() {
			continue
		}

		result = append(result, currentName)

	}

	sort.Strings(result)

	return result, nil
}
