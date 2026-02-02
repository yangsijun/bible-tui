package main

import (
	"runtime/debug"

	"github.com/yangsijun/bible-tui/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
	cmd.SetVersion(version, commit)
	cmd.Execute()
}
