# wakey

A TUI built on Charm CLI tools for managing and waking your devices using Wake-on-LAN.

## What is Wake-on-LAN?

Wake-on-LAN (WoL) is an Ethernet or Token Ring computer networking standard that allows a computer to be turned on or awakened by a network message.

The computer is woken up by sending a "magic packet" that contains the MAC address of the target computer. The magic packet is sent on the broadcast address of the network, and the target computer will turn on if the MAC address matches.

## Installation

For MacOS users, you can install `wakey` using Homebrew:

```bash
brew install wakey
```

For Windows and Linux users, you can install `wakey` using the following commands:

```bash
git clone https://github.com/jonathanruiz/wakey.git
```

## Running the application

To run the application, you can use the following command:

```bash
wakey
```

## Usage

Running the application will immediately display a list of devices that you can wake up. Within the application, you be able to add your own devices along with additional details.

You can navigate through the list using the arrow keys or VIM motions navigation and press `Enter` to wake up the selected device.

You can also press `ctrl + h` to display all the available keybindings.

## Configuration

When running `wakey` for the first time, a configuration file will be created with a list of empty devices. After the first run, `wakey` will use the configuration file to store and retrieve the devices.

The configuration file is located in your home directory at `~/.wakey_config.json`.

You can add your own devices to the configuration file by adding the following JSON object:

```json
{
  "devices": [
    {
      "DeviceName": "Device Name",
      "Description": "Description",
      "MACAddress": "00:00:00:00:00:00",
      "IPAddress": "0.0.0.0",
      "Status": "Offline"
    }
  ]
}
```

- `DeviceName` is the name of the device that you want to wake up.
- `Description` is a brief description of the device.
- `MACAddress` is the MAC address of the device.
- `IPAddress` is the IP address of the device.
- `Status` is the status of the device. This will be updated by the application. It will ping the device to determine if it is online or offline.

## Contributing

If you would like to contribute to the project, please feel free to fork the repository and submit a pull request.
