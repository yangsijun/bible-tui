package cmd

import (
	"testing"

	"github.com/sijun-dong/bible-tui/internal/db"
)

func setupTestDB(t *testing.T) *db.DB {
	t.Helper()
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })

	err = database.Migrate()
	if err != nil {
		t.Fatal(err)
	}

	vID, err := database.InsertVersion("GAE", "개역개정", "ko")
	if err != nil {
		t.Fatal(err)
	}

	bookID, err := database.InsertBook(vID, "gen", "창세기", "창", "old", 50, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = database.InsertVerse(bookID, 1, 1, "태초에 하나님이 천지를 창조하시니라", "천지 창조", false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = database.InsertVerse(bookID, 1, 2, "땅이 혼돈하고 공허하며 흑암이 깊음 위에 있고 하나님의 영은 수면 위에 운행하시니라", "", true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = database.InsertVerse(bookID, 1, 3, "하나님이 이르시되 빛이 있으라 하시니 빛이 있었고", "", false)
	if err != nil {
		t.Fatal(err)
	}

	return database
}
