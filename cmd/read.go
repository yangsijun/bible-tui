package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/sijun-dong/bible-tui/internal/bible"
	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <책이름> <장>[:<절>[-<끝절>]]",
	Short: "성경 본문 읽기",
	Long:  "지정한 성경 본문을 출력합니다. 예: bible read 창세기 1, bible read 창 1:3-5",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runRead,
}

var readVersion string

func init() {
	readCmd.Flags().StringVarP(&readVersion, "version", "v", "GAE", "성경 버전 코드")
	rootCmd.AddCommand(readCmd)
}

func runRead(cmd *cobra.Command, args []string) error {
	input := strings.Join(args, " ")
	ref, err := bible.ParseReference(input)
	if err != nil {
		return fmt.Errorf("parse reference: %w", err)
	}

	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	verses, err := database.GetVerses(readVersion, ref.BookCode, ref.Chapter)
	if err != nil {
		return fmt.Errorf("get verses: %w", err)
	}

	if len(verses) == 0 {
		return fmt.Errorf("no verses found for %s %d", ref.BookCode, ref.Chapter)
	}

	if ref.VerseStart > 0 {
		filtered := []db.Verse{}
		for _, v := range verses {
			if v.VerseNum >= ref.VerseStart {
				if ref.VerseEnd > 0 {
					if v.VerseNum <= ref.VerseEnd {
						filtered = append(filtered, v)
					}
				} else {
					if v.VerseNum == ref.VerseStart {
						filtered = append(filtered, v)
					}
				}
			}
		}
		verses = filtered
	}

	if len(verses) == 0 {
		return fmt.Errorf("no verses found in specified range")
	}

	bookName := verses[0].BookName
	chapter := verses[0].Chapter

	titleStyle := lipgloss.NewStyle().Bold(true)
	verseNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), titleStyle.Render(fmt.Sprintf("%s %d장", bookName, chapter)))
	fmt.Fprintln(cmd.OutOrStdout())

	for _, v := range verses {
		if v.SectionTitle != "" {
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintln(cmd.OutOrStdout(), titleStyle.Render(v.SectionTitle))
			fmt.Fprintln(cmd.OutOrStdout())
		}

		verseNumStr := fmt.Sprintf("%d", v.VerseNum)
		paddedNum := padLeft(verseNumStr, 3)
		coloredNum := verseNumStyle.Render(paddedNum)

		fmt.Fprintf(cmd.OutOrStdout(), "%s  %s\n", coloredNum, v.Text)
	}

	fmt.Fprintln(cmd.OutOrStdout())

	return nil
}

func padLeft(s string, width int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := strings.Repeat(" ", width-currentWidth)
	return padding + s
}
