package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestBookmarkAdd(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "add", "창", "1:1"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "책갈피 추가") {
		t.Errorf("expected '책갈피 추가' in output, got: %s", output)
	}
	if !strings.Contains(output, "창세기 1:1") {
		t.Errorf("expected '창세기 1:1' in output, got: %s", output)
	}
}

func TestBookmarkAddWithNote(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "add", "창", "1:1", "--note", "my note here"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "책갈피 추가") {
		t.Errorf("expected '책갈피 추가' in output, got: %s", output)
	}
}

func TestBookmarkList(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "add", "창", "1:1"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error adding bookmark: %v", err)
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "list"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error listing bookmarks: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "창세기 1:1") {
		t.Errorf("expected '창세기 1:1' in output, got: %s", output)
	}
	if !strings.Contains(output, "태초에 하나님이 천지를 창조하시니라") {
		t.Errorf("expected verse text in output, got: %s", output)
	}
}

func TestBookmarkRemove(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "add", "창", "1:1"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error adding bookmark: %v", err)
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "remove", "1"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error removing bookmark: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "책갈피 삭제") {
		t.Errorf("expected '책갈피 삭제' in output, got: %s", output)
	}
	if !strings.Contains(output, "#1") {
		t.Errorf("expected '#1' in output, got: %s", output)
	}
}

func TestHighlightAdd(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "add", "창", "1:1", "--color", "green"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "하이라이트 추가") {
		t.Errorf("expected '하이라이트 추가' in output, got: %s", output)
	}
	if !strings.Contains(output, "green") {
		t.Errorf("expected 'green' in output, got: %s", output)
	}
	if !strings.Contains(output, "창세기 1:1") {
		t.Errorf("expected '창세기 1:1' in output, got: %s", output)
	}
}

func TestHighlightAddDefaultColor(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	highlightColor = "yellow"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "add", "창", "1:1"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "하이라이트 추가") {
		t.Errorf("expected '하이라이트 추가' in output, got: %s", output)
	}
	if !strings.Contains(output, "yellow") {
		t.Errorf("expected 'yellow' (default color) in output, got: %s", output)
	}
}

func TestHighlightList(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "add", "창", "1:1", "--color", "blue"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error adding highlight: %v", err)
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "list"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error listing highlights: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[blue]") {
		t.Errorf("expected '[blue]' in output, got: %s", output)
	}
	if !strings.Contains(output, "창세기 1:1") {
		t.Errorf("expected '창세기 1:1' in output, got: %s", output)
	}
	if !strings.Contains(output, "태초에 하나님이 천지를 창조하시니라") {
		t.Errorf("expected verse text in output, got: %s", output)
	}
}

func TestHighlightRemove(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "add", "창", "1:1", "--color", "pink"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error adding highlight: %v", err)
	}

	buf.Reset()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "remove", "창", "1:1"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error removing highlight: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "하이라이트 삭제") {
		t.Errorf("expected '하이라이트 삭제' in output, got: %s", output)
	}
	if !strings.Contains(output, "창세기 1:1") {
		t.Errorf("expected '창세기 1:1' in output, got: %s", output)
	}
}

func TestBookmarkCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "add") {
		t.Errorf("expected 'add' subcommand in help, got: %s", output)
	}
	if !strings.Contains(output, "list") {
		t.Errorf("expected 'list' subcommand in help, got: %s", output)
	}
	if !strings.Contains(output, "remove") {
		t.Errorf("expected 'remove' subcommand in help, got: %s", output)
	}
}

func TestHighlightCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "add") {
		t.Errorf("expected 'add' subcommand in help, got: %s", output)
	}
	if !strings.Contains(output, "list") {
		t.Errorf("expected 'list' subcommand in help, got: %s", output)
	}
	if !strings.Contains(output, "remove") {
		t.Errorf("expected 'remove' subcommand in help, got: %s", output)
	}
}

func TestBookmarkAddInvalidVerse(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"bookmark", "add", "창"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for invalid reference, got nil")
	}
}

func TestHighlightAddInvalidColor(t *testing.T) {
	database := setupTestDB(t)
	testDB = database
	defer func() { testDB = nil }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"highlight", "add", "창", "1:1", "--color", "red"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for invalid color, got nil")
	}
}
