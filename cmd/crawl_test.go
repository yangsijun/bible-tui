package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sijun-dong/bible-tui/internal/db"
)

func TestCrawlCommand_DryRun(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"crawl", "--dry-run"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "스키마") && !strings.Contains(output, "완료") {
		t.Errorf("expected output to contain '스키마' or '완료', got: %s", output)
	}
}

func TestCrawlCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"crawl", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "version") {
		t.Errorf("expected output to contain 'version', got: %s", output)
	}
	if !strings.Contains(output, "book") {
		t.Errorf("expected output to contain 'book', got: %s", output)
	}
	if !strings.Contains(output, "dry-run") {
		t.Errorf("expected output to contain 'dry-run', got: %s", output)
	}
}
