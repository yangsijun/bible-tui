package main

import "github.com/yangsijun/bible-tui/cmd"

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.SetVersion(version, commit)
	cmd.Execute()
}
