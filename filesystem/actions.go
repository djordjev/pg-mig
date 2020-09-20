package filesystem

import "path/filepath"

// CreateMigrationFile - creates a new file in path directory
func (fs *Filesystem) CreateMigrationFile(name string, location string) error {
	filename := filepath.Join(location, name)

	_, err := fs.Fs.Create(filename)
	if err != nil {
		return err
	}

	return nil
}
