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

//func TestLoadFiles(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	time1, _ := time.Parse(time.RFC3339, "2020-09-20T15:04:05Z")
//	time2, _ := time.Parse(time.RFC3339, "2020-09-20T16:04:05Z")
//	time3, _ := time.Parse(time.RFC3339, "2020-09-20T17:04:05Z")
//	time4, _ := time.Parse(time.RFC3339, "2020-09-20T18:04:05Z")
//
//	afero.WriteFile(fs, fmt.Sprintf("mig_%d_something_up.sql", time1.Unix()), []byte("first"), os.ModePerm)
//	afero.WriteFile(fs, fmt.Sprintf("mig_%d_up.sql", time2.Unix()), []byte("second"), os.ModePerm)
//	afero.WriteFile(fs, fmt.Sprintf("mig_%d_zzzz_up.sql", time3.Unix()), []byte("third"), os.ModePerm)
//	afero.WriteFile(fs, fmt.Sprintf("mig_%d_up.sql", time1.Unix()), []byte("fourth"), os.ModePerm)
//	afero.WriteFile(fs, "invalid_file_format", []byte("invalid"), os.ModePerm)
//
//	rexp := regexp.MustCompile("^mig_([0-9]+).*_up.sql$")
//
//	table := []struct {
//		from     time.Time
//		to       time.Time
//		contents []string
//		ids      []int64
//	}{
//		{
//			from:    time1.Sub(),
//			to:      time4,
//			pattern: regexp.MustCompile("^mig_([0-9]+).*_up.sql$"),
//		},
//	}
//}
