![](https://raw.githubusercontent.com/SuperFromND/iguana/master/res/logo.svg)
---
***Iguana*** (**I**KEMEN **G**O **U**tility for **A**nnotating **N**onspecified **A**ttacks) is a tool for generating [movelist.dat](https://github.com/ikemen-engine/Ikemen-GO/wiki/Miscellaneous-Info#movelists) files for [I.K.E.M.E.N GO and M.U.G.E.N](https://github.com/ikemen-engine/Ikemen-GO) characters, using only their standard command definitions (.cmd) file. It is the successor/continuation of my JavaScript-based [IKEMEN GO Command List Generator](https://superfromnd.gitlab.io/ikemen-cmdlist/), written in Go and released as a standalone executable.

Iguana is a work-in-progress as of this writing (August 2023), and is not yet feature-complete. Despite this, it is considered functional enough for a public release.

## Usage
Iguana is a command line tool that requires at least one parameter, `-i`, to specify a command file to load.
```bash
$ iguana.exe -i path/to/file.cmd
```

When given a directory, Iguana can bulk-process all of the .cmd files it can find within that directory.
```bash
$ iguana.exe -i path/to/folder
```

By default, Iguana will output a file named `movelist.dat` in the same directory as the input file is located in. This location cannot currently be changed, but you *can* change its filename using the `-o` parameter.

Additional options and parameters can be found by running Iguana with either no arguments at all or with the `-h` argument.

## Features
Iguana is a work in progress, so right now it's missing a lot of functionality from its web-based incarnation. Several features are planned to be implemented in the future, including:
- [ ] Formatting options
  - [x] Use Special Button Labels
  - [x] Compress Motions
  - [ ] Label Moves with Header Comments
  - [x] Remove one-button commands
  - [ ] Annotate palette-specific moves
  - [ ] Annotate air-state moves
- [ ] Customizable header colors
- [ ] Button remapping support
- [ ] Automatic `[Movelist]` patching of character.def files
- [x] Bulk processing of entire `chars` folders

## Building
Iguana requires at least Go version 1.16 to compile. Aside from that, you shouldn't have too much trouble building.
```bash
$ git clone https://github.com/SuperFromND/iguana.git
$ cd iguana
$ make
```
The resulting executable will be placed in the `bin` directory.

## Licensing
Iguana, as well as its logo, is [available under the MIT License.](https://raw.githubusercontent.com/SuperFromND/iguana/main/LICENSE) <3