# lrcsnc
Gets the currently playing song's synced lyrics from https://lrclib.net/ (if there are) and displays them in sync with player's position.

lrcsnc is primarily designed for bars like [waybar](https://github.com/Alexays/Waybar).

https://github.com/user-attachments/assets/209ddfdd-0c2a-4ce6-a213-b9796a154c28

## Build
```
git clone https://github.com/Endg4meZer0/lrcsnc.git
cd lrcsnc
go get
go build
```
Should do the trick.

## Usage
```
lrcsnc [OPTION]
```
Get more info on on available options with `lrcsnc -help` or on [wiki](https://github.com/Endg4meZer0/lrcsnc/wiki/Available-options).

## TODO
- [ ] A TUI implementation for... showing off rices, maybe
- [ ] Additional functionality for bars (maybe by implementing IPC)
- [ ] More configuration options?
- [ ] There is definitely always more!

## Known issues
- Spotify: if you leave songs on autoplay without using previous or next buttons, lyrics may desync a lot. It's an internal issue of Spotify's position data desyncing from the song's actual position and is not related to lrcsnc, MPRIS or even D-Bus. This may get fixed by itself after 7-8 seconds if the songs are not from the same album, but you can also fix it manually (e.g. pause & play, seek on the position bar, or even toggling the shuffle or repeat mode works).

## Not a known issue or you have an enhancement suggestion?
Please, make an issue so I can fix it, suggest a workaround or add a new feature!

## A song was not found on LrcLib?
Consider adding the lyrics for it! LrcLib is a great open-source lyrics provider service that has its own easy-to-use [app](https://github.com/tranxuanthang/lrcget) to download or upload lyrics. Once the lyrics are uploaded, lrcsnc should be able to pick them up on the next play of the song if the cached version of said song's lyrics is outdated/not found. If the cached version exists, you may delete it using the existing flags (check wiki for more info on that).
