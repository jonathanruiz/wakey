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
		name, _ := cmd.Flags().GetString("name")
		devicesFlag, _ := cmd.Flags().GetString("devices")

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

func init() {
	groupCmd.Flags().StringP("name", "n", "", "name of the group to wake")

	groupCreateCmd.Flags().StringP("name", "n", "", "group name")
	groupCreateCmd.Flags().StringP("devices", "d", "", "comma-separated list of device names to add to the group")

	groupCmd.AddCommand(groupCreateCmd)
}
