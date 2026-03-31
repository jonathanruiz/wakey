package main

import (
	"fmt"
	"os"
	"wakey/internal/common/wol"
	"wakey/internal/config"

	"github.com/spf13/cobra"
)

var wakeCmd = &cobra.Command{
	Use:   "wake",
	Short: "Wake a device or group",
	RunE: func(cmd *cobra.Command, args []string) error {
		device, _ := cmd.Flags().GetString("device")
		group, _ := cmd.Flags().GetString("group")

		if device == "" && group == "" {
			return cmd.Help()
		}

		cfg := config.ReadConfig()

		if device != "" {
			if err := wakeDeviceCLI(device, cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		if group != "" {
			if err := wakeGroupCLI(group, cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		return nil
	},
}

func init() {
	wakeCmd.Flags().StringP("device", "d", "", "name of the device to wake")
	wakeCmd.Flags().StringP("group", "g", "", "name of the group to wake")
}

func wakeDeviceCLI(name string, cfg config.Config) error {
	for _, d := range cfg.Devices {
		if d.DeviceName == name {
			if err := wol.WakeDevice(d.MacAddress); err != nil {
				return fmt.Errorf("error waking device %q: %w", name, err)
			}
			fmt.Printf("Waking device %q (%s)\n", d.DeviceName, d.MacAddress)
			return nil
		}
	}
	return fmt.Errorf("device not found: %q", name)
}

func wakeGroupCLI(name string, cfg config.Config) error {
	for _, g := range cfg.Groups {
		if g.GroupName == name {
			deviceMACMap := make(map[string]string)
			for _, d := range cfg.Devices {
				deviceMACMap[d.ID] = d.MacAddress
			}

			var macs []string
			for _, id := range g.Devices {
				if mac, ok := deviceMACMap[id]; ok {
					macs = append(macs, mac)
				}
			}

			if err := wol.WakeGroup(macs); err != nil {
				return fmt.Errorf("error waking group %q: %w", name, err)
			}
			fmt.Printf("Waking group %q (%d device(s))\n", g.GroupName, len(macs))
			return nil
		}
	}
	return fmt.Errorf("group not found: %q", name)
}
