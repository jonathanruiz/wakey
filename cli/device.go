package cli

import (
	"fmt"
	"os"
	"regexp"
	"wakey/internal/common/wol"
	"wakey/internal/config"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Wake a device by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the device name from the command flags and validate it.
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return cmd.Help()
		}

		// Read the application configuration and find the device with the specified name. If found, wake the device using its MAC address.
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

var deviceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new device",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the device details from the command flags and validate them.
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		mac, _ := cmd.Flags().GetString("mac")
		ip, _ := cmd.Flags().GetString("ip")

		// Validate that all required fields are provided.
		if name == "" {
			return fmt.Errorf("device name is required (-n)")
		}
		if description == "" {
			return fmt.Errorf("description is required (-d)")
		}
		if mac == "" {
			return fmt.Errorf("MAC address is required (-m)")
		}
		if ip == "" {
			return fmt.Errorf("IP address is required (-i)")
		}

		// Validate the MAC address format (basic validation).
		macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
		if !macRegex.MatchString(mac) {
			return fmt.Errorf("invalid MAC address: %q", mac)
		}

		// Validate the IP address format (basic validation).
		var a, b, c, d int
		if _, err := fmt.Sscanf(ip, "%d.%d.%d.%d", &a, &b, &c, &d); err != nil {
			return fmt.Errorf("invalid IP address: %q", ip)
		}

		cfg := config.ReadConfig()
		cfg.Devices = append(cfg.Devices, config.Device{
			ID:          uuid.NewString(),
			DeviceName:  name,
			Description: description,
			MacAddress:  mac,
			IPAddress:   ip,
			State:       "Offline",
		})
		config.WriteConfig(cfg)

		fmt.Printf("Device %q created\n", name)
		return nil
	},
}

var deviceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if a device is online",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the device name from the command flags and validate it.
		name, _ := cmd.Flags().GetString("name")

		if name == "" {
			return fmt.Errorf("device name is required (-n)")
		}

		cfg := config.ReadConfig()
		for _, d := range cfg.Devices {
			if d.DeviceName == name {
				if wol.IsOnline(d.IPAddress) {
					fmt.Printf("%s is online\n", d.DeviceName)
				} else {
					fmt.Printf("%s is offline\n", d.DeviceName)
				}
				return nil
			}
		}

		fmt.Fprintf(os.Stderr, "Device not found: %q\n", name)
		os.Exit(1)
		return nil
	},
}

func init() {
	// wakey device -n <devicename>
	deviceCmd.Flags().StringP("name", "n", "", "name of the device to wake")

	// wakey device create -n <devicename> -d <description> -m <MAC> -i <IP>
	deviceCreateCmd.Flags().StringP("name", "n", "", "device name")
	deviceCreateCmd.Flags().StringP("description", "d", "", "device description")
	deviceCreateCmd.Flags().StringP("mac", "m", "", "MAC address (e.g. AA:BB:CC:DD:EE:FF)")
	deviceCreateCmd.Flags().StringP("ip", "i", "", "IP address (e.g. 192.168.1.100)")

	// wakey device status -n <devicename>
	deviceStatusCmd.Flags().StringP("name", "n", "", "name of the device to check")

	deviceCmd.AddCommand(deviceCreateCmd)
	deviceCmd.AddCommand(deviceStatusCmd)
}
