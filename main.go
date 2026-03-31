package main

import (
	"fmt"
	"os"

	"wakey/internal"
	"wakey/internal/common/status"
	"wakey/internal/config"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "wakey",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(internal.InitialModel())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(groupCmd)
}

func main() {
	status.Message = config.CreateConfig()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
