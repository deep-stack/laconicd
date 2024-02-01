package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"git.vdb.to/cerc-io/laconic2d/app"
	"git.vdb.to/cerc-io/laconic2d/app/params"
	"git.vdb.to/cerc-io/laconic2d/cmd/laconic2d/cmd"
)

func main() {
	params.SetAddressPrefixes()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
