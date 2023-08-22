![](https://raw.githubusercontent.com/SuperFromND/iguana/master/res/logo.svg)
---
***Iguana*** (**I**KEMEN **G**O **U**tility for **A**nnotating **N**onspecified **A**ttacks) is a tool for generating [movelist.dat](https://github.com/ikemen-engine/Ikemen-GO/wiki/Miscellaneous-Info#movelists) files for [I.K.E.M.E.N GO and M.U.G.E.N](https://github.com/ikemen-engine/Ikemen-GO) characters, using only their standard command definitions (.cmd) file. It is the successor/continuation of my JavaScript-based [IKEMEN GO Command List Generator](https://superfromnd.gitlab.io/ikemen-cmdlist/), written in Go and released as a standalone executable.

Iguana is a work-in-progress as of this writing (August 2023), and is not yet feature-complete or even feature-congruent with the aforementioned IKEMEN GO Command List Generator. Despite this, it is nevertheless functional enough to be deemed worthy of a public release.

## Usage
Iguana is a command line tool that requires at least one parameter, `-i`, to specify a command file to load. Additional options and parameters can be found by running Iguana by itself or with the `-h` argument.
```bash
$ iguana.exe -i path/to/file.cmd
```

By default, Iguana will output a file named `movelist.dat` in the same directory as the input file is located in. This location cannot currently be changed, but you *can* change its filename using the `-o` parameter.

## Features
Iguana is a work in progress, so right now it isn't capable of much outside of its basic functionality. Several features are planned to be implemented in the future, including:
- [ ] Formatting options
  - [ ] Use Special Button Labels
  - [ ] Compress Motions
  - [ ] Label Moves with Header Comments
  - [x] Remove one-button commands
  - [ ] Annotate palette-specific moves
  - [ ] Annotate air-state moves
- [ ] Customizable header colors
- [ ] Button remapping support
- [ ] Automatic `[Movelist]` patching of character.def files
- [ ] Bulk processing of entire `chars` folders

## Building
Iguana requires at least Go version 1.13 to compile, due to [the INI library](https://github.com/go-ini/ini) requiring it. Aside from that, you shouldn't have too much trouble building as it's otherwise standard Go.
```
$ git clone https://github.com/SuperFromND/iguana.git
$ cd iguana
$ make
```

## Licensing
Iguana, as well as its logo, is [available under the MIT License.](https://raw.githubusercontent.com/SuperFromND/iguana/main/LICENSE) <3