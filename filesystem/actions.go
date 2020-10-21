package filesystem

import (
	"fmt"
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
		return fmt.Errorf("filesystem error: unable to create migration file %w", err)
	}

	return nil
}

// DeleteMigrationFiles - removed given migration files from FS
func (fs *ImplFilesystem) DeleteMigrationFiles(list MigrationFileList) error {
	for _, file := range list {
		up, down, err := fs.getMigrationsForTimestamp(file.Timestamp)
		if err != nil {
			return err
		}

		err = fs.Fs.Remove(up)
		if err != nil {
			return err
		}

		err = fs.Fs.Remove(down)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *ImplFilesystem) getMigrationsForTimestamp(ts int64) (up string, down string, e error) {
	config, err := fs.LoadConfig()
	if err != nil {
		e = err
		return
	}

	files, err := afero.ReadDir(fs.Fs, config.Path)
	if err != nil {
		e = fmt.Errorf("filesystem error: unable to read config file directory %w", err)
		return
	}

	upPattern := regexp.MustCompile(fmt.Sprintf("^mig_%d*_up.sql$", ts))
	downPattern := regexp.MustCompile(fmt.Sprintf("^mig_%d*_down.sql$", ts))

	for _, file := range files {
		if upPattern.FindString(file.Name()) != "" {
			up = filepath.Join(config.Path, file.Name())
		}

		if downPattern.FindString(file.Name()) != "" {
			down = filepath.Join(config.Path, file.Name())
		}

		if up != "" && down != "" {
			return
		}
	}

	return "", "", fmt.Errorf("filesystem error: unable to find migrations for %d", ts)
}

// ReadMigrationContent for specified migration file reads content as a string
func (fs *ImplFilesystem) ReadMigrationContent(file MigrationFile, direction Direction, config Config) (string, error) {
	path := file.GetFileName(config, direction)

	content, err := afero.ReadFile(fs.Fs, path)
	if err != nil {
		return "", fmt.Errorf("filesystem error: unable to read migration file content %w", err)
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
		return nil, fmt.Errorf("filesystem error: unable to read config file directory %w", err)
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
		return fmt.Errorf("filesystem error: given filename does not match regex %+q", submatches)
	}

	fullName := submatches[0]
	strTs := submatches[1]
	intTs, err := strconv.ParseInt(strTs, 10, 64)
	if err != nil {
		return fmt.Errorf("filesystem error: invalid migration file name %+q %w", submatches, err)
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
