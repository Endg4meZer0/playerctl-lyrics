# playerctl-lyrics
Uses `playerctl metadata` output to get song data and makes a GET request to https://lrclib.net/ to get the song's synced lyrics. These lyrics are then printed to stdout in sync with `playerctl position`.

This small thingie is primarily designed for bars like [waybar](https://github.com/Alexays/Waybar).

https://github.com/user-attachments/assets/b3107a9a-eabf-4d76-a6da-ac82fc911569

## TODO
- [ ] Better error handling
- [ ] Configuration
- [ ] Flag usage
- [ ] Some simple caching system
- [ ] There is always more!
