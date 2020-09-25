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

func (fs *Filesystem) GetFileTimestamps(from time.Time, to time.Time) (MigrationFileList, error) {
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
		result = append(result, v)
	}

	sort.Sort(result)

	return result, nil
}

func (fs *Filesystem) storeMigrationFileInMap(resultMap map[int64]MigrationFile, submatches []string, isUp bool) error {
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

//func (fs *Filesystem) GetFiles(from time.Time, to time.Time, direction Direction) ([]MigrationFile, error) {
//	config, err := fs.LoadConfig()
//	if err != nil {
//		return nil, err
//	}
//
//	files, err := afero.ReadDir(fs.Fs, config.Path)
//	if err != nil {
//		return nil, err
//	}
//
//	if direction == DirectionUp {
//		pattern := regexp.MustCompile("^mig_([0-9]+).*_up.sql$")
//		files, err := fs.loadFiles(from, to, pattern, files)
//		if err != nil {
//			return nil, err
//		}
//		sort.Sort(files)
//
//		return files, nil
//	}
//
//	if direction == DirectionDown {
//		pattern := regexp.MustCompile("^mig_([0-9]+).*_down.sql$")
//		files, err := fs.loadFiles(from, to, pattern, files)
//		if err != nil {
//			return nil, err
//		}
//
//		sort.Sort(sort.Reverse(files))
//
//		return files, nil
//	}
//
//	return nil, errors.New("invalid direction")
//}
//
//func (fs *Filesystem) loadFiles(from time.Time, to time.Time, pattern *regexp.Regexp, files []os.FileInfo) (FileList, error) {
//	result := make([]MigrationFile, 0, 10)
//
//	for i := 0; i < len(files); i++ {
//		current := files[i]
//		currentName := current.Name()
//
//		matched := pattern.FindString(currentName)
//		if matched == "" {
//			continue
//		}
//
//		ts := pattern.FindStringSubmatch(currentName)
//		if ts == nil || len(ts) < 2 {
//			return nil, InvalidFileFormat
//		}
//
//		strTs := ts[1]
//		intTs, err := strconv.ParseInt(strTs, 10, 64)
//		if err != nil {
//			return nil, err
//		}
//
//		if intTs < from.Unix() || intTs > to.Unix() {
//			continue
//		}
//
//		fileContent, err := afero.ReadFile(fs.Fs, currentName)
//		if err != nil {
//			return nil, err
//		}
//
//		result = append(result, MigrationFile{Timestamp: intTs, Content: string(fileContent)})
//
//	}
//
//	return result, nil
//}
