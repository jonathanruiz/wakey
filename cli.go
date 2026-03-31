package main

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

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Wake a group by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return cmd.Help()
		}

		cfg := config.ReadConfig()
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
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("Waking group %q (%d device(s))\n", g.GroupName, len(macs))
				return nil
			}
		}

		fmt.Fprintf(os.Stderr, "Group not found: %q\n", name)
		os.Exit(1)
		return nil
	},
}

func init() {
	deviceCmd.Flags().StringP("name", "n", "", "name of the device to wake")
	groupCmd.Flags().StringP("name", "n", "", "name of the group to wake")
}
