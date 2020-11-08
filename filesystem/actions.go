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

	err = fs.Fs.Chmod(filename, 0666)
	if err != nil {
		return fmt.Errorf("filesystem error: unable to set permissions on file %w", err)
	}

	return nil
}

// DeleteMigrationFiles - removed given migration files from FS
func (fs *ImplFilesystem) deleteMigrationFiles(list MigrationFileList) error {
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

// Squash squashes files from given list into one up migration and one down migration
func (fs *ImplFilesystem) Squash(files MigrationFileList) (err error) {
	if len(files) == 0 {
		err = fmt.Errorf("filesystem error: No files to squash")
		return
	}

	config, err := fs.LoadConfig()
	if err != nil {
		return err
	}

	up := make([]string, 0, len(files))
	down := make([]string, 0, len(files))

	for _, file := range files {
		upContent, err := fs.ReadMigrationContent(file, DirectionUp, config)
		if err != nil {
			return err
		}

		upComment := fmt.Sprintf("-- migration %d UP\n", file.Timestamp)
		up = append(up, fmt.Sprintf("%s%s%s", upComment, upContent, "\n"))

		downContent, err := fs.ReadMigrationContent(file, DirectionDown, config)
		if err != nil {
			return err
		}

		downComment := fmt.Sprintf("-- migration %d DOWN\n", file.Timestamp)
		down = append(down, fmt.Sprintf("%s%s%s", downComment, downContent, "\n"))
	}

	lastTs := files[len(files)-1]
	upName := fmt.Sprintf("mig_%d_%s_up.sql", lastTs.Timestamp, "squashed")
	err = fs.writeFile(up, upName, config)
	if err != nil {
		return
	}

	// Reverse down migrations
	for i := 0; i < len(down)/2; i++ {
		j := len(down) - i - 1
		down[i], down[j] = down[j], down[i]
	}

	downName := fmt.Sprintf("mig_%d_%s_down.sql", lastTs.Timestamp, "squashed")
	err = fs.writeFile(down, downName, config)
	if err != nil {
		return
	}

	err = fs.deleteMigrationFiles(files)
	if err != nil {
		return
	}

	return
}

func (fs *ImplFilesystem) writeFile(migrations []string, name string, config Config) (err error) {
	filename := filepath.Join(config.Path, name)

	file, err := fs.Fs.Create(filename)
	if err != nil {
		err = fmt.Errorf("filesystem error: unable to create migration file for squash %w", err)
		return
	}

	err = fs.Fs.Chmod(filename, 0666)
	if err != nil {
		return fmt.Errorf("filesystem error: unable to set permissions on file %w", err)
	}


	defer file.Close()

	for _, mig := range migrations {
		file.WriteString(mig)
	}

	return
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
