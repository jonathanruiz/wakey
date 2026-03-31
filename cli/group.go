package cli

import (
	"fmt"
	"os"
	"strings"
	"wakey/internal/common/wol"
	"wakey/internal/config"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Wake a group by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the group name from the command flags and validate it.
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return cmd.Help()
		}

		// Read the application configuration and find the group with the specified name. If found, wake all devices in the group using their MAC addresses.
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

var groupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new group",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the group name and device names from the command flags and validate them.
		name, _ := cmd.Flags().GetString("name")
		devicesFlag, _ := cmd.Flags().GetString("devices")

		// Validate that the group name is provided.
		if name == "" {
			return fmt.Errorf("group name is required (-n)")
		}

		cfg := config.ReadConfig()

		// Resolve device names to IDs
		deviceNameToID := make(map[string]string)
		for _, d := range cfg.Devices {
			deviceNameToID[d.DeviceName] = d.ID
		}

		// Parse the comma-separated device names and look up their IDs. If any device name is not found, return an error.
		var deviceIDs []string
		if devicesFlag != "" {
			for _, deviceName := range strings.Split(devicesFlag, ",") {
				deviceName = strings.TrimSpace(deviceName)
				id, ok := deviceNameToID[deviceName]
				if !ok {
					return fmt.Errorf("device not found: %q", deviceName)
				}
				deviceIDs = append(deviceIDs, id)
			}
		}

		cfg.Groups = append(cfg.Groups, config.Group{
			ID:        uuid.NewString(),
			GroupName: name,
			Devices:   deviceIDs,
		})
		config.WriteConfig(cfg)

		fmt.Printf("Group %q created\n", name)
		return nil
	},
}

var groupStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if devices in a group are online",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the group name from the command flags and validate it.
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("group name is required (-n)")
		}

		cfg := config.ReadConfig()
		for _, g := range cfg.Groups {
			if g.GroupName == name {
				deviceMap := make(map[string]config.Device)

				// Build a map of device ID to device struct for quick lookup.
				for _, d := range cfg.Devices {
					deviceMap[d.ID] = d
				}

				// Check the status of each device in the group and print whether it is online or offline.
				for _, id := range g.Devices {
					d, ok := deviceMap[id]
					if !ok {
						continue
					}
					if wol.IsOnline(d.IPAddress) {
						fmt.Printf("%s is online\n", d.DeviceName)
					} else {
						fmt.Printf("%s is offline\n", d.DeviceName)
					}
				}
				return nil
			}
		}

		fmt.Fprintf(os.Stderr, "Group not found: %q\n", name)
		os.Exit(1)
		return nil
	},
}

func init() {
	// wakey group -n <groupname>
	groupCmd.Flags().StringP("name", "n", "", "name of the group to wake")

	// wakey group create -n <groupname> -d <device1,device2,...>
	groupCreateCmd.Flags().StringP("name", "n", "", "group name")
	groupCreateCmd.Flags().StringP("devices", "d", "", "comma-separated list of device names to add to the group")

	// wakey group status -n <groupname>
	groupStatusCmd.Flags().StringP("name", "n", "", "name of the group to check")

	groupCmd.AddCommand(groupCreateCmd)
	groupCmd.AddCommand(groupStatusCmd)
}
