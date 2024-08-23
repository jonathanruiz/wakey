# Where should we write the GIF?
Output test/demo.gif

# Set up a 1200x600 terminal with 46px font.
Set FontSize 16
Set Width 1200
Set Height 500

# Type a command in the terminal.
Type "wakey"

# Pause for dramatic effect...
Sleep 500ms

# Run the command by pressing enter.
Enter

# Admire the output for a bit.
Sleep 5s

# Create a new device
Type "c"

# Enter device name
Type "demo-03"

# Go to next input
Tab

# Enter device description
Type@50ms "This is a server using for this demo!"

# Go to next input
Tab

# Enter device MAC Address
Type "a1:b2:c3:d4:e5:f6"

# Go to next input
Tab

# Enter device IP Address
Type "192.168.0.15"

# Go to next input
Tab

# Wait
Sleep 2s

# Submit device
Enter

# Wait
Sleep 1s

# Move down
Type "j"

# Wait
Sleep 500ms

# Move down
Type "j"

# Wait
Sleep 500ms

# Move down
Type "j"

# Wait
Sleep 1s

# Wake device
Enter

Sleep 5s

Hide
Type "d"
Type "y"
Escape
