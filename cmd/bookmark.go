package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sijun-dong/bible-tui/internal/bible"
	"github.com/spf13/cobra"
)

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "책갈피 관리",
	Long:  "성경 책갈피를 추가, 삭제, 목록 조회합니다.",
}

var highlightCmd = &cobra.Command{
	Use:   "highlight",
	Short: "하이라이트 관리",
	Long:  "성경 하이라이트를 추가, 삭제, 목록 조회합니다.",
}

// Bookmark subcommands
var bookmarkAddCmd = &cobra.Command{
	Use:   "add <참조> [--note \"메모\"]",
	Short: "책갈피 추가",
	Long:  "성경 구절에 책갈피를 추가합니다. 예: bible bookmark add 창 1:1",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runBookmarkAdd,
}

var bookmarkListCmd = &cobra.Command{
	Use:   "list",
	Short: "책갈피 목록",
	Long:  "저장된 책갈피 목록을 조회합니다.",
	RunE:  runBookmarkList,
}

var bookmarkRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "책갈피 삭제",
	Long:  "책갈피 ID로 책갈피를 삭제합니다.",
	Args:  cobra.ExactArgs(1),
	RunE:  runBookmarkRemove,
}

// Highlight subcommands
var highlightAddCmd = &cobra.Command{
	Use:   "add <참조> [--color yellow]",
	Short: "하이라이트 추가",
	Long:  "성경 구절에 하이라이트를 추가합니다. 예: bible highlight add 창 1:1 --color green",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runHighlightAdd,
}

var highlightListCmd = &cobra.Command{
	Use:   "list",
	Short: "하이라이트 목록",
	Long:  "저장된 하이라이트 목록을 조회합니다.",
	RunE:  runHighlightList,
}

var highlightRemoveCmd = &cobra.Command{
	Use:   "remove <참조>",
	Short: "하이라이트 삭제",
	Long:  "성경 구절의 하이라이트를 삭제합니다. 예: bible highlight remove 창 1:1",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runHighlightRemove,
}

// Flags
var (
	bookmarkNote      string
	bookmarkVersion   string
	bookmarkLimit     int
	highlightColor    string
	highlightVersion  string
	highlightLimit    int
)

func init() {
	// Bookmark flags
	bookmarkAddCmd.Flags().StringVar(&bookmarkNote, "note", "", "책갈피 메모")
	bookmarkAddCmd.Flags().StringVarP(&bookmarkVersion, "version", "v", "GAE", "성경 버전 코드")
	bookmarkListCmd.Flags().IntVar(&bookmarkLimit, "limit", 20, "조회할 책갈피 개수")
	bookmarkRemoveCmd.Flags().StringVarP(&bookmarkVersion, "version", "v", "GAE", "성경 버전 코드")

	// Highlight flags
	highlightAddCmd.Flags().StringVar(&highlightColor, "color", "yellow", "하이라이트 색상 (yellow, green, blue, pink, purple)")
	highlightAddCmd.Flags().StringVarP(&highlightVersion, "version", "v", "GAE", "성경 버전 코드")
	highlightListCmd.Flags().IntVar(&highlightLimit, "limit", 20, "조회할 하이라이트 개수")
	highlightRemoveCmd.Flags().StringVarP(&highlightVersion, "version", "v", "GAE", "성경 버전 코드")

	// Register subcommands
	bookmarkCmd.AddCommand(bookmarkAddCmd, bookmarkListCmd, bookmarkRemoveCmd)
	highlightCmd.AddCommand(highlightAddCmd, highlightListCmd, highlightRemoveCmd)

	// Register parent commands
	rootCmd.AddCommand(bookmarkCmd)
	rootCmd.AddCommand(highlightCmd)
}

// resolveVerseID resolves a reference to a verse ID and label
func resolveVerseID(versionCode string, args []string) (int64, string, error) {
	ref, err := bible.ParseReference(strings.Join(args, " "))
	if err != nil {
		return 0, "", err
	}

	if ref.VerseStart == 0 {
		return 0, "", fmt.Errorf("specific verse required (e.g., 창 1:1)")
	}

	database, err := getDB()
	if err != nil {
		return 0, "", fmt.Errorf("open database: %w", err)
	}

	verses, err := database.GetVerses(versionCode, ref.BookCode, ref.Chapter)
	if err != nil {
		return 0, "", err
	}

	for _, v := range verses {
		if v.VerseNum == ref.VerseStart {
			label := fmt.Sprintf("%s %d:%d", v.BookName, v.Chapter, v.VerseNum)
			return v.ID, label, nil
		}
	}

	return 0, "", fmt.Errorf("verse not found: %s %d:%d", ref.BookCode, ref.Chapter, ref.VerseStart)
}

// Bookmark command implementations
func runBookmarkAdd(cmd *cobra.Command, args []string) error {
	verseID, label, err := resolveVerseID(bookmarkVersion, args)
	if err != nil {
		return err
	}

	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	_, err = database.AddBookmark(verseID, bookmarkNote)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "책갈피 추가: %s\n", label)
	return nil
}

func runBookmarkList(cmd *cobra.Command, args []string) error {
	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	bookmarks, err := database.ListBookmarks(bookmarkLimit, 0)
	if err != nil {
		return err
	}

	if len(bookmarks) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "저장된 책갈피가 없습니다.")
		return nil
	}

	for _, bm := range bookmarks {
		fmt.Fprintf(cmd.OutOrStdout(), "[ID:%d] %s %d:%d — %s\n",
			bm.ID, bm.BookName, bm.Chapter, bm.VerseNum, bm.VerseText)
		if bm.Note != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "  메모: %s\n", bm.Note)
		}
	}

	return nil
}

func runBookmarkRemove(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid bookmark ID: %w", err)
	}

	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	err = database.RemoveBookmark(id)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "책갈피 삭제: #%d\n", id)
	return nil
}

// Highlight command implementations
func runHighlightAdd(cmd *cobra.Command, args []string) error {
	// Validate color
	validColors := map[string]bool{
		"yellow": true,
		"green":  true,
		"blue":   true,
		"pink":   true,
		"purple": true,
	}
	if !validColors[highlightColor] {
		return fmt.Errorf("invalid color: %s (valid: yellow, green, blue, pink, purple)", highlightColor)
	}

	verseID, label, err := resolveVerseID(highlightVersion, args)
	if err != nil {
		return err
	}

	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	err = database.AddHighlight(verseID, highlightColor)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "하이라이트 추가: %s (%s)\n", label, highlightColor)
	return nil
}

func runHighlightList(cmd *cobra.Command, args []string) error {
	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	highlights, err := database.ListHighlights(highlightLimit, 0)
	if err != nil {
		return err
	}

	if len(highlights) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "저장된 하이라이트가 없습니다.")
		return nil
	}

	for _, h := range highlights {
		fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s %d:%d — %s\n",
			h.Color, h.BookName, h.Chapter, h.VerseNum, h.VerseText)
	}

	return nil
}

func runHighlightRemove(cmd *cobra.Command, args []string) error {
	verseID, label, err := resolveVerseID(highlightVersion, args)
	if err != nil {
		return err
	}

	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	err = database.RemoveHighlight(verseID)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "하이라이트 삭제: %s\n", label)
	return nil
}
