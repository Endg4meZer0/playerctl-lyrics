# playerctl-lyrics
Uses `playerctl metadata` output to get song data and makes a GET request to https://lrclib.net/ to get the song's synced lyrics (if there are). These lyrics are then printed to stdout in sync with `playerctl position`.

This small thingie is primarily designed for bars like [waybar](https://github.com/Alexays/Waybar).

https://github.com/user-attachments/assets/b3107a9a-eabf-4d76-a6da-ac82fc911569

## Important!
Players that use seconds as position data like `cmus` will cause a lot of lyrics to mismatch with their accurate timings, potentionally skipping some at all. It *will* work, but nonetheless it is advised to use other players that have more precise format for position data. For example, Spotify works good, since it uses milliseconds in position data.

## Build
`go build` inside the directory should do the trick.

## TODO
- [ ] Better error handling
- [ ] Configuration
- [ ] Some simple caching system
- [ ] Better handling of players with seconds as position data
- [ ] Flag usage
- [ ] There is always more!

## Known issues
- If you leave songs on autoplay in Spotify without using previous or next buttons, lyrics may desync a lot. It's an internal issue of Spotify position data desyncing from the song's actual position and is not related to playerctl-lyrics. This can be fixed just by pausing playback and continuing it again.

## A song was not found on LrcLib?
Consider adding the lyrics for it! LrcLib is a great open-source lyrics provider service that has its own easy-to-use [app](https://github.com/tranxuanthang/lrcget) to download or upload lyrics. Once the lyrics are uploaded, `playerctl-lyrics` should be able to pick them up on the next play of the song.