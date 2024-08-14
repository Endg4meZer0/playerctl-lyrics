# playerctl-lyrics
Uses `playerctl metadata` output to get song data and makes a GET request to https://lrclib.net/ to get the song's synced lyrics (if there are). These lyrics are then printed to stdout in sync with `playerctl position`.

This small thingie is primarily designed for bars like [waybar](https://github.com/Alexays/Waybar).

https://github.com/user-attachments/assets/b3107a9a-eabf-4d76-a6da-ac82fc91156

## Important!
Players that use seconds as position data like `cmus` will cause a lot of lyrics to mismatch with their accurate timings, potentionally skipping some at all. It *will* work, but nonetheless it is advised to use other players that have more precise format for position data. For example, Spotify works good, since it uses milliseconds in position data.

## Build
```
git clone https://github.com/Endg4meZer0/playerctl-lyrics.git
cd playerctl-lyrics
go get
go build
```
Should do the trick.

## Usage
```
playerctl-lyrics [OPTION]
```
Get more info on on available options with `playerctl-lyrics -help` or on wiki.

## TODO
- [x] ~~Better error handling~~
- [x] ~~Caching system~~
- [x] ~~Options handling~~
- [x] ~~Configuration~~
- [ ] Better handling of players with seconds as position data
- [ ] There is always more!

## Known issues
- If you leave songs on autoplay in Spotify without using previous or next buttons, lyrics may desync a lot. It's an internal issue of Spotify position data desyncing from the song's actual position and is not related to playerctl-lyrics. This sometimes gets fixed by itself during the playback, but it can also be fixed manually by pausing playback and continuing it again.

## A song was not found on LrcLib?
Consider adding the lyrics for it! LrcLib is a great open-source lyrics provider service that has its own easy-to-use [app](https://github.com/tranxuanthang/lrcget) to download or upload lyrics. Once the lyrics are uploaded, `playerctl-lyrics` should be able to pick them up on the next play of the song.
