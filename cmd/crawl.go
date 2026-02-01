package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yangsijun/bible-tui/internal/crawler"
)

var (
	crawlVersion string
	crawlBook    string
	crawlDryRun  bool
	crawlReset   bool
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "성경 데이터 크롤링",
	Long:  "대한성서공회 웹사이트에서 성경 데이터를 크롤링합니다.",
	RunE:  runCrawl,
}

func init() {
	crawlCmd.Flags().StringVar(&crawlVersion, "version", "GAE", "version code")
	crawlCmd.Flags().StringVar(&crawlBook, "book", "", "specific book code to crawl (empty = all)")
	crawlCmd.Flags().BoolVar(&crawlDryRun, "dry-run", false, "only create DB schema, don't crawl")
	crawlCmd.Flags().BoolVar(&crawlReset, "reset", false, "delete crawled data and re-crawl")
	rootCmd.AddCommand(crawlCmd)
}

func runCrawl(cmd *cobra.Command, args []string) error {
	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer database.Close()

	// Always run migration
	if err := database.Migrate(); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if crawlDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "DB 스키마 생성 완료\n")
		return nil
	}

	if crawlReset {
		deleted, err := database.ResetCrawlData(crawlVersion, crawlBook)
		if err != nil {
			return fmt.Errorf("reset crawl data: %w", err)
		}
		target := crawlVersion
		if crawlBook != "" {
			target += "/" + crawlBook
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s 크롤링 데이터 삭제 완료 (구절 %d개)\n", target, deleted)
	}

	c := crawler.New(
		database,
		crawler.WithVersionCode(crawlVersion),
		crawler.WithOnProgress(func(bookName string, chapter, totalChapters int) {
			fmt.Fprintf(cmd.OutOrStdout(), "[%d/%d] %s %d장 크롤링 완료\n", chapter, totalChapters, bookName, chapter)
		}),
	)

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Crawl based on flags
	if crawlBook != "" {
		if err := c.CrawlBook(ctx, crawlBook); err != nil {
			return fmt.Errorf("crawl book: %w", err)
		}
	} else {
		if err := c.CrawlAll(ctx); err != nil {
			return fmt.Errorf("crawl all: %w", err)
		}
	}

	return nil
}
