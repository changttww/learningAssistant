package config

import "testing"

func TestLoadConfigDefaultsToMySQL(t *testing.T) {
	t.Setenv("DB_DRIVER", "")
	t.Setenv("DB_SQLITE_PATH", "")

	LoadConfig()

	if AppConfig.Database.Driver != "mysql" {
		t.Fatalf("expected default database driver mysql, got %q", AppConfig.Database.Driver)
	}
	if AppConfig.Database.SQLitePath != "learning_assistant.db" {
		t.Fatalf("expected default sqlite path, got %q", AppConfig.Database.SQLitePath)
	}
}

func TestLoadConfigReadsSQLiteSettings(t *testing.T) {
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_SQLITE_PATH", "tmp/dev.db")

	LoadConfig()

	if AppConfig.Database.Driver != "sqlite" {
		t.Fatalf("expected sqlite database driver, got %q", AppConfig.Database.Driver)
	}
	if AppConfig.Database.SQLitePath != "tmp/dev.db" {
		t.Fatalf("expected configured sqlite path, got %q", AppConfig.Database.SQLitePath)
	}
}
