# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - 202X-XX-XX
### Added
- Terminal User Interface using [bubbletea](https://github.com/charmbracelet/bubbletea)
- Some simple unit tests like cache and romanization
### Changed
- A-a-a-and another **ginormous** refactor, learning and making use of moduling and SOLID principles.
- Changed configuration correspondingly (you will learn more on wiki once the version releases); planning to move configuration format from JSON to KDL.
### Removed
- Playerctl support is now removed in favor of direct MPRIS/D-Bus handling.
- Terminal output in one line is now removed since a proper TUI is now availabe.

## [[0.2.1](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.2.1)] - 2024-08-29
### Added
- A command-line option `-o` to redirect the output to a set file.
- ~~A command-line option to display lyrics in one line~~ **is deprecated**
- A configuration option to offset the lyrics by set seconds by @Endg4meZer0 in [#9](https://github.com/Endg4meZer0/playerctl-lyrics/pull/9)
### Changed
- More refactoring: `cmus` and other players that report position in integer seconds are now fully supported.
- Cache system is reverted back to JSON instead of LRC files to allow more additional data to be stored ([#10](https://github.com/Endg4meZer0/playerctl-lyrics/pull/10))
### Fixed
- Instrumental lyrics overlapped actual lyrics in some cases ([#11](https://github.com/Endg4meZer0/playerctl-lyrics/pull/11))

## [[0.2.0](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.2.0)] - 2024-08-24
### Changed
- A big concept rewrite happened to allow players like `cmus` that report position in integer seconds work on par with others.
- A rename of `doCacheLyrics` configuration option to `enabled`

## [[0.1.1](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.1.1)] - 2024-08-21
### Added
- A configuration option to control the format of repeated lyrics multiplier.
### Fixed
- Fixed a panic if there is no space after a timestamp.
- Fixed a panic when romanization of Japanese kanji failed and fell down to Chinese characters 

## [[0.1.0](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.1.0)] - 2024-08-15
### Added
- Initial unstable release of playerctl-lyrics.
- Display lyrics for currently playing song.
- Support for multiple music players using `playerctl`.
- Automatic lyric fetching from `lrclib`.
- Configuration file for custom settings.
- Romanization for several asian languages.
- Caching system to significantly reduce traffic usage.