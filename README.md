# playerctl-lyrics
Uses `playerctl metadata` output to get song data and makes a GET request to https://lrclib.net/ to get the song's synced lyrics (if there are). These lyrics are then printed to stdout in sync with `playerctl position`.

This small thingie is primarily designed for bars like [waybar](https://github.com/Alexays/Waybar).

https://github.com/user-attachments/assets/b3107a9a-eabf-4d76-a6da-ac82fc911569

## Important!
Players with integer-based position data like `cmus` will cause a lot of lyrics to mismatch with their accurate timings, potentionally skipping some at all. It *will* work, but nonetheless it is adviced to use other players that have float-based position data. For example, `spotify` works perfectly.

## TODO
- [ ] Better error handling
- [ ] Configuration
- [ ] Flag usage
- [ ] Some simple caching system
- [ ] There is always more!
