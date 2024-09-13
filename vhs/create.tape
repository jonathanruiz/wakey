# Where should we write the GIF?
Output create.gif

# Set up a 1200x600 terminal with 46px font.
Set FontSize 16
Set Width 1200
Set Height 500

Hide
# Create a new device
Type "wakey"
Sleep 500ms
Enter
Sleep 5s
Type "c"

Show

# Enter device name
Sleep 1s
Type "demo-03"
Sleep 500ms

# Go to and enter device MAC Address
Enter
Sleep 500ms
Enter
Sleep 500ms
Type "a1:b2:c3:d4:e5:f"

# Enter device IP Address
Enter
Type "10.1.34.1"

# Try to submit device
Enter
Sleep 2s
Enter

# Navigate to Mac Address
Sleep 3s
Enter
Sleep 1s
Enter
Sleep 2s
Type "6"
Sleep 2s

# Navigate to Submit
Enter
Sleep 500ms
Enter
Sleep 500ms
Enter

# Navigate to Description
Sleep 2s
Enter
Type@50ms "I forgot to add the description!"
Sleep 1s

# Navigate to Submit
Enter @ 500ms 4

# Edit device
Sleep 1s
Type "j"
Sleep 500ms
Type "j"
Sleep 1s
Type "e"
Sleep 1s

# Get to description and replace it
Enter
Sleep 1s
Backspace 32
Sleep 1s
Type@50ms "Actually, this is a server!"

# Get to Submit button and submit
Enter @ 1s 3
Sleep 1s
Enter

# Wait 
Sleep 5s

Hide
Type "d"
Type "y"
Escape
