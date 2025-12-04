package main

import (
	"github.com/analog-substance/util/cli/build_info"
	"github.com/analog-substance/util/cli/updater/cobra_updater"
	"github.com/defektive/gnas/pkg/cmd"
)

var version = "v0.0.0"
var commit = "replace"

func main() {
	versionInfo := build_info.InitLoadedVersion(version, commit)
	cmd.RootCmd.Version = versionInfo.String()
	cobra_updater.AddToRootCmd(cmd.RootCmd, versionInfo)

	cmd.Execute()
}
