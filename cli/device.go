package cli

import (
	"fmt"
	"os"
	"wakey/internal/common/wol"
	"wakey/internal/config"

	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Wake a device by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return cmd.Help()
		}

		cfg := config.ReadConfig()
		for _, d := range cfg.Devices {
			if d.DeviceName == name {
				if err := wol.WakeDevice(d.MacAddress); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("Waking device %q (%s)\n", d.DeviceName, d.MacAddress)
				return nil
			}
		}

		fmt.Fprintf(os.Stderr, "Device not found: %q\n", name)
		os.Exit(1)
		return nil
	},
}

func init() {
	deviceCmd.Flags().StringP("name", "n", "", "name of the device to wake")
}
