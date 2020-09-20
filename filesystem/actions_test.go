package filesystem

import (
	"github.com/spf13/afero"
	"testing"
)

func TestCreateFileWithLocation(t *testing.T) {
	fs := afero.NewMemMapFs()
	fsystem := Filesystem{Fs: fs}

	err := fsystem.CreateMigrationFile("test_name.sql", "./test/")
	if err != nil {
		t.Logf("Unable to create migration file %v", err)
		t.Fail()
	}

	exists, err := afero.Exists(fs, "test/test_name.sql")
	if err != nil {
		t.Logf("Error when creating a file %v", err)
		t.Fail()
	}

	if !exists {
		t.Log("Created file does not exist on FS")
		t.Fail()
	}

	exists, err = afero.Exists(fs, "test_name.sql")
	if err != nil {
		t.Logf("Error when checking created file %v", err)
		t.Fail()
	}

	if exists {
		t.Log("Created file ignored location")
		t.Fail()
	}
}

func TestCreateFileNoLocation(t *testing.T) {
	fs := afero.NewMemMapFs()
	fsystem := Filesystem{Fs: fs}

	err := fsystem.CreateMigrationFile("no_loc_name.sql", "")
	if err != nil {
		t.Logf("Unable to create migration file %v", err)
		t.Fail()
	}

	exists, err := afero.Exists(fs, "no_loc_name.sql")
	if err != nil {
		t.Logf("Error when creating a file %v", err)
		t.Fail()
	}

	if !exists {
		t.Log("Created file does not exist on FS")
		t.Fail()
	}
}
