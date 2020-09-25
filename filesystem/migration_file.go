package filesystem

import "path/filepath"

type MigrationFile struct {
	Timestamp int64
	Down      string
	Up        string
}

// GetFileName - returns full path to migration file
func (mf MigrationFile) GetFileName(config Config, direction Direction) string {
	path := config.Path

	if direction == DirectionUp {
		return filepath.Join(path, mf.Up)
	}

	return filepath.Join(path, mf.Down)
}

type MigrationFileList []MigrationFile

func (m MigrationFileList) Len() int           { return len(m) }
func (m MigrationFileList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MigrationFileList) Less(i, j int) bool { return m[i].Timestamp < m[j].Timestamp }
