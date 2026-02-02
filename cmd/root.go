package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bible",
	Short: "성경 TUI/CLI 프로그램",
	Long:  "성경을 터미널에서 읽고 검색하는 프로그램입니다.",
}

// SetVersion sets the version info from ldflags
func SetVersion(version, commit string) {
	rootCmd.Version = version
	if commit != "none" && commit != "" {
		rootCmd.Version = fmt.Sprintf("%s (%s)", version, commit)
	}
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
