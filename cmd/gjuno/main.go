package main

import (
	junocmd "github.com/forbole/juno/v4/cmd"
	initcmd "github.com/forbole/juno/v4/cmd/init"
	migratecmd "github.com/forbole/juno/v4/cmd/migrate"
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	startcmd "github.com/forbole/juno/v4/cmd/start"

	parsecmd "github.com/hjcore/gjuno/cmd/parse"
	gotabitdb "github.com/hjcore/gjuno/database"
	"github.com/hjcore/gjuno/x"
)

func main() {
	// Setup the config
	parseCfg := parsecmdtypes.NewConfig().
		WithRegistrar(x.NewModulesRegistrar()).
		WithDBBuilder(gotabitdb.Builder)

	cfg := junocmd.NewConfig("gjuno").
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := junocmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		junocmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		startcmd.NewStartCmd(cfg.GetParseConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(cfg.GetName(), cfg.GetParseConfig()),
	)

	executor := junocmd.PrepareRootCmd(cfg.GetName(), rootCmd)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
