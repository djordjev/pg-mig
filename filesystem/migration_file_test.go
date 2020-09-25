package filesystem

import "testing"

func TestMigrationGetFileName(t *testing.T) {
	mf := MigrationFile{
		Timestamp: 123456,
		Up:        "mig_123456_something_up.sql",
		Down:      "mig_123456_something_down.sql",
	}

	config := Config{Path: "/demo/path"}

	upMigration := mf.GetFileName(config, DirectionUp)
	if upMigration != "/demo/path/mig_123456_something_up.sql" {
		t.Logf("Invalid return value: %s", upMigration)
		t.Fail()
	}

	downMigration := mf.GetFileName(config, DirectionDown)
	if downMigration != "/demo/path/mig_123456_something_down.sql" {
		t.Logf("Invalid return value: %s", downMigration)
		t.Fail()
	}
}
