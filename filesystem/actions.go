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

// CreateMigrationFile - creates a new file in path directory
func (fs *ImplFilesystem) CreateMigrationFile(name string, location string) error {
	filename := filepath.Join(location, name)

	_, err := fs.Fs.Create(filename)
	if err != nil {
		return err
	}

	return nil
}

// ReadMigrationContent for specified migration file reads content as a string
func (fs *ImplFilesystem) ReadMigrationContent(file MigrationFile, direction Direction, config Config) (string, error) {
	path := file.GetFileName(config, direction)

	content, err := afero.ReadFile(fs.Fs, path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetFileTimestamps - gets the list of migrations that are between two arguments.
// Returned list does not include file that has exactly same timestamp as `from` arg.
// Returned list includes file that has exactly same timestamp as `to` arg.
func (fs *ImplFilesystem) GetFileTimestamps(from time.Time, to time.Time) (MigrationFileList, error) {
	config, err := fs.LoadConfig()
	if err != nil {
		return nil, err
	}

	files, err := afero.ReadDir(fs.Fs, config.Path)
	if err != nil {
		return nil, err
	}

	upPattern := regexp.MustCompile("^mig_([0-9]+).*_up.sql$")
	downPattern := regexp.MustCompile("^mig_([0-9]+).*_down.sql$")

	resultMap := make(map[int64]MigrationFile)

	for i := 0; i < len(files); i++ {
		current := files[i]
		currentName := current.Name()

		upSubmatches := upPattern.FindStringSubmatch(currentName)
		downSubmatches := downPattern.FindStringSubmatch(currentName)

		if upSubmatches != nil && len(upSubmatches) >= 2 {
			err = fs.storeMigrationFileInMap(resultMap, upSubmatches, true)
			if err != nil {
				return nil, err
			}
		}

		if downSubmatches != nil && len(downSubmatches) >= 2 {
			err = fs.storeMigrationFileInMap(resultMap, downSubmatches, false)
			if err != nil {
				return nil, err
			}
		}
	}

	result := make(MigrationFileList, 0, len(resultMap))

	for _, v := range resultMap {
		if v.Timestamp > from.Unix() && v.Timestamp <= to.Unix() {
			result = append(result, v)
		}
	}

	sort.Sort(result)

	return result, nil
}

func (fs *ImplFilesystem) storeMigrationFileInMap(resultMap map[int64]MigrationFile, submatches []string, isUp bool) error {
	if len(submatches) < 2 {
		return errors.New("given filename does not match regex")
	}

	fullName := submatches[0]
	strTs := submatches[1]
	intTs, err := strconv.ParseInt(strTs, 10, 64)
	if err != nil {
		return err
	}

	entry, ok := resultMap[intTs]
	if !ok {
		migration := MigrationFile{Timestamp: intTs}
		if isUp {
			migration.Up = fullName
		} else {
			migration.Down = fullName
		}
		resultMap[intTs] = migration
	} else {
		if isUp {
			entry.Up = fullName
		} else {
			entry.Down = fullName
		}
		resultMap[intTs] = entry
	}

	return nil
}
