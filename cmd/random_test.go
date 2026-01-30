package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRandomCommand(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"random"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output")
	}

	if !strings.Contains(output, "창세기") {
		t.Errorf("expected output to contain book name '창세기', got: %s", output)
	}

	if !strings.Contains(output, "1:") {
		t.Errorf("expected output to contain chapter:verse format '1:', got: %s", output)
	}
}
