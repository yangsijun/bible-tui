package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestSearchCommand(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"search", "하나님"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "창세기") {
		t.Errorf("expected output to contain '창세기', got: %s", output)
	}
	if !strings.Contains(output, "하나님") {
		t.Errorf("expected output to contain '하나님', got: %s", output)
	}
	if !strings.Contains(output, "검색 결과") {
		t.Errorf("expected output to contain '검색 결과', got: %s", output)
	}
}

func TestSearchCommand_NoResults(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"search", "존재하지않는단어"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "검색 결과가 없습니다") {
		t.Errorf("expected output to contain '검색 결과가 없습니다', got: %s", output)
	}
	if !strings.Contains(output, "존재하지않는단어") {
		t.Errorf("expected output to contain '존재하지않는단어', got: %s", output)
	}
}

func TestSearchCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"search", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "limit") {
		t.Errorf("expected output to contain 'limit', got: %s", output)
	}
	if !strings.Contains(output, "version") {
		t.Errorf("expected output to contain 'version', got: %s", output)
	}
}
