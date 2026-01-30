package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "랜덤 성경 구절",
	Long:  "랜덤으로 성경 구절 하나를 출력합니다.",
	RunE:  runRandom,
}

var randomVersion string

func init() {
	randomCmd.Flags().StringVarP(&randomVersion, "version", "v", "GAE", "성경 버전 코드")
	rootCmd.AddCommand(randomCmd)
}

func runRandom(cmd *cobra.Command, args []string) error {
	database, err := getDB()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	verse, err := database.GetRandomVerse(randomVersion)
	if err != nil {
		return fmt.Errorf("get random verse: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s %d:%d — %s\n",
		verse.BookName, verse.Chapter, verse.VerseNum, verse.Text)

	return nil
}
