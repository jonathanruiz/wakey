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
