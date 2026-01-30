package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sijun-dong/bible-tui/internal/bible"
	"github.com/sijun-dong/bible-tui/internal/db"
)

func loadFixture(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../../testdata/genesis_1.html")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return string(data)
}

func setupTestDB(t *testing.T) *db.DB {
	t.Helper()
	d, err := db.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	if err := d.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func seedVersionAndBooks(t *testing.T, d *db.DB) int64 {
	t.Helper()
	versionID, err := d.InsertVersion("GAE", "개역개정", "ko")
	if err != nil {
		t.Fatalf("InsertVersion: %v", err)
	}
	for i, b := range bible.AllBooks() {
		_, err := d.InsertBook(versionID, b.Code, b.NameKo, b.AbbrevKo, b.Testament, b.ChapterCount, i)
		if err != nil {
			t.Fatalf("InsertBook %s: %v", b.Code, err)
		}
	}
	return versionID
}

func mockServer(t *testing.T, fixture string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fixture))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestCrawlBook_Mock(t *testing.T) {
	fixture := loadFixture(t)
	d := setupTestDB(t)
	srv := mockServer(t, fixture)

	seedVersionAndBooks(t, d)

	c := New(d,
		WithBaseURL(srv.URL),
		WithRateLimit(1000),
	)

	ctx := context.Background()
	// Obadiah has 1 chapter — fastest possible end-to-end test
	if err := c.CrawlBook(ctx, "oba"); err != nil {
		t.Fatalf("CrawlBook(oba): %v", err)
	}

	verses, err := d.GetVerses("GAE", "oba", 1)
	if err != nil {
		t.Fatalf("GetVerses: %v", err)
	}
	if len(verses) != 31 {
		t.Errorf("expected 31 verses (genesis_1 fixture), got %d", len(verses))
	}
	if verses[0].Text == "" {
		t.Error("first verse text is empty")
	}

	status, err := d.GetCrawlStatus("GAE", "oba", 1)
	if err != nil {
		t.Fatalf("GetCrawlStatus: %v", err)
	}
	if status != "done" {
		t.Errorf("expected crawl status 'done', got %q", status)
	}

	doneCount, err := d.CountCrawlDone("GAE")
	if err != nil {
		t.Fatalf("CountCrawlDone: %v", err)
	}
	if doneCount != 1 {
		t.Errorf("expected 1 done, got %d", doneCount)
	}
}

func TestResumability(t *testing.T) {
	fixture := loadFixture(t)
	d := setupTestDB(t)

	var requestCount atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fixture))
	}))
	t.Cleanup(srv.Close)

	seedVersionAndBooks(t, d)

	// Haggai has 2 chapters — mark chapter 1 as done before crawling
	if err := d.SetCrawlStatus("GAE", "hag", 1, "done", 10, ""); err != nil {
		t.Fatalf("SetCrawlStatus: %v", err)
	}

	c := New(d,
		WithBaseURL(srv.URL),
		WithRateLimit(1000),
	)

	ctx := context.Background()
	if err := c.CrawlBook(ctx, "hag"); err != nil {
		t.Fatalf("CrawlBook(hag): %v", err)
	}

	// Only chapter 2 should have been fetched
	if got := requestCount.Load(); got != 1 {
		t.Errorf("expected 1 request (ch2 only), got %d", got)
	}

	status, err := d.GetCrawlStatus("GAE", "hag", 2)
	if err != nil {
		t.Fatalf("GetCrawlStatus ch2: %v", err)
	}
	if status != "done" {
		t.Errorf("expected 'done' for ch2, got %q", status)
	}
}

func TestRateLimiting(t *testing.T) {
	fixture := loadFixture(t)
	d := setupTestDB(t)
	srv := mockServer(t, fixture)

	seedVersionAndBooks(t, d)

	// 10 req/sec → 100ms between requests; 3 chapters of Nahum → ≥200ms
	c := New(d,
		WithBaseURL(srv.URL),
		WithRateLimit(10),
	)

	ctx := context.Background()
	start := time.Now()
	if err := c.CrawlBook(ctx, "nam"); err != nil {
		t.Fatalf("CrawlBook(nam): %v", err)
	}
	elapsed := time.Since(start)

	if elapsed < 180*time.Millisecond {
		t.Errorf("expected ≥180ms for 3 chapters at 10rps, got %v", elapsed)
	}
}

func TestValidation(t *testing.T) {
	fixture := loadFixture(t)
	d := setupTestDB(t)
	srv := mockServer(t, fixture)

	seedVersionAndBooks(t, d)

	c := New(d,
		WithBaseURL(srv.URL),
		WithRateLimit(1000),
	)

	ctx := context.Background()

	// Crawl all chapters of all 66 books so validation has full data
	// That's too slow — instead, test Validate failure mode
	err := c.Validate(ctx)
	if err == nil {
		t.Fatal("expected Validate to fail when no verses exist")
	}

	// Now seed one book fully (oba = 1 chapter) and test partial validation fails
	if err := c.CrawlBook(ctx, "oba"); err != nil {
		t.Fatalf("CrawlBook(oba): %v", err)
	}

	// Validation still fails because not all 66 books have verses
	err = c.Validate(ctx)
	if err == nil {
		t.Fatal("expected Validate to fail with incomplete data")
	}
}

func TestFetchChapter_ErrorHandling(t *testing.T) {
	d := setupTestDB(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	t.Cleanup(srv.Close)

	seedVersionAndBooks(t, d)

	c := New(d,
		WithBaseURL(srv.URL),
		WithRateLimit(1000),
	)

	ctx := context.Background()
	// Crawl a 1-chapter book — should not crash, should record error status
	err := c.CrawlBook(ctx, "oba")
	if err != nil {
		t.Fatalf("CrawlBook should not return error on chapter failure, got: %v", err)
	}

	status, err := d.GetCrawlStatus("GAE", "oba", 1)
	if err != nil {
		t.Fatalf("GetCrawlStatus: %v", err)
	}
	if status != "error" {
		t.Errorf("expected 'error' status, got %q", status)
	}
}

func TestContextCancellation(t *testing.T) {
	fixture := loadFixture(t)
	d := setupTestDB(t)
	srv := mockServer(t, fixture)

	seedVersionAndBooks(t, d)

	c := New(d,
		WithBaseURL(srv.URL),
		WithRateLimit(1000),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := c.CrawlBook(ctx, "gen")
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestWithOptions(t *testing.T) {
	d := setupTestDB(t)

	custom := &http.Client{Timeout: 5 * time.Second}
	var called bool
	progressFn := func(bookName string, chapter, total int) { called = true }

	c := New(d,
		WithBaseURL("http://example.com"),
		WithVersionCode("HAN"),
		WithVersionName("개역한글"),
		WithHTTPClient(custom),
		WithRateLimit(100),
		WithOnProgress(progressFn),
	)

	if c.baseURL != "http://example.com" {
		t.Errorf("baseURL: got %q", c.baseURL)
	}
	if c.versionCode != "HAN" {
		t.Errorf("versionCode: got %q", c.versionCode)
	}
	if c.versionName != "개역한글" {
		t.Errorf("versionName: got %q", c.versionName)
	}
	if c.client != custom {
		t.Error("client not set")
	}
	if c.onProgress == nil {
		t.Error("onProgress not set")
	}

	// Trigger the progress callback to verify it works
	c.onProgress("test", 1, 1)
	if !called {
		t.Error("progress callback not invoked")
	}
}
