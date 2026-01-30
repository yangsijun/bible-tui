package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/time/rate"

	"github.com/sijun-dong/bible-tui/internal/bible"
	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/parser"
)

type Crawler struct {
	db          *db.DB
	client      *http.Client
	limiter     *rate.Limiter
	baseURL     string
	versionCode string
	versionName string
	onProgress  func(bookName string, chapter, totalChapters int)
}

type Option func(*Crawler)

func WithBaseURL(url string) Option {
	return func(c *Crawler) { c.baseURL = url }
}

func WithVersionCode(code string) Option {
	return func(c *Crawler) { c.versionCode = code }
}

func WithVersionName(name string) Option {
	return func(c *Crawler) { c.versionName = name }
}

func WithHTTPClient(cl *http.Client) Option {
	return func(c *Crawler) { c.client = cl }
}

func WithRateLimit(rps float64) Option {
	return func(c *Crawler) { c.limiter = rate.NewLimiter(rate.Limit(rps), 1) }
}

func WithOnProgress(fn func(bookName string, chapter, totalChapters int)) Option {
	return func(c *Crawler) { c.onProgress = fn }
}

func New(database *db.DB, opts ...Option) *Crawler {
	c := &Crawler{
		db:          database,
		client:      &http.Client{Timeout: 30 * time.Second},
		limiter:     rate.NewLimiter(rate.Limit(0.5), 1),
		baseURL:     "https://www.bskorea.or.kr/bible/korbibReadpage.php",
		versionCode: "GAE",
		versionName: "개역개정",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Crawler) CrawlAll(ctx context.Context) error {
	versionID, err := c.db.InsertVersion(c.versionCode, c.versionName, "ko")
	if err != nil {
		return fmt.Errorf("insert version: %w", err)
	}

	books := bible.AllBooks()
	for i, b := range books {
		_, err := c.db.InsertBook(versionID, b.Code, b.NameKo, b.AbbrevKo, b.Testament, b.ChapterCount, i)
		if err != nil {
			return fmt.Errorf("insert book %s: %w", b.Code, err)
		}
	}

	for _, b := range books {
		if err := c.crawlBookChapters(ctx, b.Code, b.NameKo, b.ChapterCount); err != nil {
			return err
		}
	}

	return c.Validate(ctx)
}

func (c *Crawler) CrawlBook(ctx context.Context, bookCode string) error {
	info, ok := bible.GetBookByCode(bookCode)
	if !ok {
		return fmt.Errorf("unknown book code: %s", bookCode)
	}
	return c.crawlBookChapters(ctx, info.Code, info.NameKo, info.ChapterCount)
}

func (c *Crawler) crawlBookChapters(ctx context.Context, bookCode, bookName string, chapterCount int) error {
	for ch := 1; ch <= chapterCount; ch++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		status, err := c.db.GetCrawlStatus(c.versionCode, bookCode, ch)
		if err != nil {
			return fmt.Errorf("get crawl status %s ch%d: %w", bookCode, ch, err)
		}
		if status == "done" {
			if c.onProgress != nil {
				c.onProgress(bookName, ch, chapterCount)
			}
			continue
		}

		if err := c.crawlChapter(ctx, bookCode, ch); err != nil {
			_ = c.db.SetCrawlStatus(c.versionCode, bookCode, ch, "error", 0, err.Error())
			continue
		}

		if c.onProgress != nil {
			c.onProgress(bookName, ch, chapterCount)
		}
	}
	return nil
}

func (c *Crawler) crawlChapter(ctx context.Context, bookCode string, chapter int) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}

	htmlBody, err := c.fetchChapter(ctx, bookCode, chapter)
	if err != nil {
		return fmt.Errorf("fetch %s ch%d: %w", bookCode, chapter, err)
	}

	parsed, err := parser.ParseChapterHTML(htmlBody)
	if err != nil {
		return fmt.Errorf("parse %s ch%d: %w", bookCode, chapter, err)
	}

	book, err := c.db.GetBookByCode(c.versionCode, bookCode)
	if err != nil {
		return fmt.Errorf("get book %s: %w", bookCode, err)
	}

	for _, v := range parsed.Verses {
		verseID, err := c.db.InsertVerse(book.ID, chapter, v.Number, v.Text, v.SectionTitle, len(v.Footnotes) > 0)
		if err != nil {
			return fmt.Errorf("insert verse %s %d:%d: %w", bookCode, chapter, v.Number, err)
		}
		for _, fn := range v.Footnotes {
			if err := c.db.InsertFootnote(verseID, fn.Marker, fn.Content); err != nil {
				return fmt.Errorf("insert footnote %s %d:%d: %w", bookCode, chapter, v.Number, err)
			}
		}
	}

	return c.db.SetCrawlStatus(c.versionCode, bookCode, chapter, "done", len(parsed.Verses), "")
}

func (c *Crawler) Validate(ctx context.Context) error {
	books := bible.AllBooks()
	for _, b := range books {
		_, err := c.db.GetBookByCode(c.versionCode, b.Code)
		if err != nil {
			return fmt.Errorf("validate: book %s not found: %w", b.Code, err)
		}

		for ch := 1; ch <= b.ChapterCount; ch++ {
			if err := ctx.Err(); err != nil {
				return err
			}
			verses, err := c.db.GetVerses(c.versionCode, b.Code, ch)
			if err != nil {
				return fmt.Errorf("validate: get verses %s ch%d: %w", b.Code, ch, err)
			}
			if len(verses) == 0 {
				return fmt.Errorf("validate: no verses for %s chapter %d", b.Code, ch)
			}
			for _, v := range verses {
				if v.Text == "" {
					return fmt.Errorf("validate: empty verse text %s %d:%d", b.Code, ch, v.VerseNum)
				}
			}
		}
	}
	return nil
}

func (c *Crawler) fetchChapter(ctx context.Context, bookCode string, chapter int) (string, error) {
	url := fmt.Sprintf("%s?version=%s&book=%s&chap=%d", c.baseURL, c.versionCode, bookCode, chapter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "BibleTUI/1.0 (Personal; non-commercial)")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status %d", resp.StatusCode)
	}

	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("charset reader: %w", err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}
