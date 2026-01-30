package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <검색어>",
	Short: "성경 검색",
	Long:  "성경 본문에서 검색어를 검색합니다.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

var (
	searchVersion string
	searchLimit   int
)

func init() {
	searchCmd.Flags().StringVarP(&searchVersion, "version", "v", "GAE", "성경 버전 코드")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "검색 결과 최대 개수")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	results, err := database.SearchVerses(searchVersion, query, searchLimit)
	if err != nil {
		return fmt.Errorf("search verses: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "검색 결과가 없습니다: \"%s\"\n", query)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\"%s\" 검색 결과 (%d건):\n", query, len(results))

	highlightStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))

	for i, result := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s %d:%d\n", i+1, result.Verse.BookName, result.Verse.Chapter, result.Verse.VerseNum)

		snippet := result.Snippet
		if snippet == "" {
			snippet = result.Verse.Text
		}

		highlightedSnippet := highlightSearchTerm(snippet, query, highlightStyle)
		fmt.Fprintf(cmd.OutOrStdout(), "    %s\n", highlightedSnippet)
	}

	return nil
}

func highlightSearchTerm(text, query string, style lipgloss.Style) string {
	if query == "" {
		return text
	}

	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	var result strings.Builder
	pos := 0

	for {
		idx := strings.Index(lowerText[pos:], lowerQuery)
		if idx == -1 {
			result.WriteString(text[pos:])
			break
		}

		result.WriteString(text[pos : pos+idx])

		matchText := text[pos+idx : pos+idx+len(query)]
		result.WriteString(style.Render(matchText))

		pos += idx + len(query)
	}

	return result.String()
}
