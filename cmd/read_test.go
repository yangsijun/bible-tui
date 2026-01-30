package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadCommand_Chapter(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"read", "창세기", "1"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "태초에") {
		t.Errorf("expected output to contain '태초에', got: %s", output)
	}
	if !strings.Contains(output, "창세기 1장") {
		t.Errorf("expected output to contain '창세기 1장', got: %s", output)
	}
}

func TestReadCommand_VerseRange(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"read", "창세기", "1:2-3"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "태초에") {
		t.Errorf("expected output to NOT contain verse 1 '태초에', got: %s", output)
	}
	if !strings.Contains(output, "땅이 혼돈하고") {
		t.Errorf("expected output to contain verse 2 '땅이 혼돈하고', got: %s", output)
	}
	if !strings.Contains(output, "빛이 있으라") {
		t.Errorf("expected output to contain verse 3 '빛이 있으라', got: %s", output)
	}
}

func TestReadCommand_InvalidBook(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"read", "없는책", "1"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid book, got nil")
	}

	if !strings.Contains(err.Error(), "unknown book") {
		t.Errorf("expected error to contain 'unknown book', got: %v", err)
	}
}
