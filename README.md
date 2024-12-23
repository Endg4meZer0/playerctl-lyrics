# App is undergoing a huge refactor and feature additions before v0.3.0, check other branches for more info

# playerctl-lyrics
Gets data using playerctl and the currently playing song's synced lyrics from https://lrclib.net/ (if there are). Lyrics are displayed in sync with player's position.

This small thingie is primarily designed for bars like [waybar](https://github.com/Alexays/Waybar).

https://github.com/user-attachments/assets/209ddfdd-0c2a-4ce6-a213-b9796a154c28

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
Get more info on on available options with `playerctl-lyrics -help` or on [wiki](https://github.com/Endg4meZer0/playerctl-lyrics/wiki/Available-options).

## TODO
- [x] ~~Caching system~~
- [x] ~~Flags/options handling~~
- [x] ~~Configuration~~
- [x] ~~Better handling of players with seconds as position data~~
- [x] ~~An ability to redirect the lyrics output to a file (same concept, one lyric at a time in sync)~~
- [ ] More different configuration options?
- [ ] There is always more!

## Known issues
- Spotify: if you leave songs on autoplay without using previous or next buttons, lyrics may desync a lot. It's an internal issue of Spotify's reported position data desyncing from the song's actual position and is not related to playerctl-lyrics. This sometimes gets fixed by itself during the playback, but it can also be fixed manually by pausing playback and continuing it again or seeking to anywhere on the position bar.

## Not a known issue or you have an enhancement suggestion?
Please, make an issue so I can fix it, suggest a workaround or add a new feature!

## A song was not found on LrcLib?
Consider adding the lyrics for it! LrcLib is a great open-source lyrics provider service that has its own easy-to-use [app](https://github.com/tranxuanthang/lrcget) to download or upload lyrics. Once the lyrics are uploaded, playerctl-lyrics should be able to pick them up on the next play of the song if the cached version of said song's lyrics is outdated/not found. If the cached version exists, you may delete it using the existing flags (check wiki for more info on that).
