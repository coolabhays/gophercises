package cmd

import "github.com/spf13/cobra"

// TODO:
// add subcommand "rm"
// add subcommand "complete"
// checkout "https://github.com/gophercises/task" to know more

var rootCmd = &cobra.Command {
	Use: "task",
	Short: "task is a CLI task manager",
	Long: `A bare minimum CLI task manager which is being used to
	show different stuffs related to cobra and boltDB`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
