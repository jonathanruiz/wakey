# Where should we write the GIF?
Output test/demo.gif

# Set up a 1200x600 terminal with 46px font.
Set FontSize 16
Set Width 1200
Set Height 500

Type "wakey"

Sleep 500ms

Enter

Sleep 5s

# Create a new device
Type "c"

# Enter device name
Type "demo-03"

# Go to next input
Enter

# Enter device description
Type@50ms "This is a mini pc for the demo."

# Go to next input
Enter

# Enter device MAC Address
Type "a1:b2:c3:d4:e5:f6"

# Go to next input
Enter

# Enter device IP Address
Type "10.1.34.1"

# Go to next input
Enter

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

# Wait
Sleep 5s

# Refresh devices
Type "r"

# Refresh devices
Sleep 5s

# Move down
Type "j"

# Wait
Sleep 500ms

# Move down
Type "j"

# Wait
Sleep 1s

# Edit device
Type "e"

# Wait
Sleep 1s

# Get to description
Enter

# Wait
Sleep 1s

# Erase description
Backspace 31 

# Wait
Sleep 1s

# Enter new description
Type@50ms "Actually, this is a server!"

# Get to Submit button
Enter @ 1s 3

# Wait
Sleep 1s

# Submit changes
Enter

# Wait 
Sleep 5s

Hide
Type "d"
Type "y"
Escape
