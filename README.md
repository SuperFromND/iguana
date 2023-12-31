![](https://raw.githubusercontent.com/SuperFromND/iguana/master/res/logo.svg)
---
***Iguana*** (**I**kemen **G**O **U**tility for **A**nnotating **N**onspecified **A**ttacks) is a tool for generating [movelist.dat](https://github.com/Ikemen-engine/Ikemen-GO/wiki/Miscellaneous-Info#movelists) files for [Ikemen GO and M.U.G.E.N](https://github.com/Ikemen-engine/Ikemen-GO) characters, using only their standard command definitions (.cmd) file. It is the successor/continuation of my JavaScript-based [Ikemen GO Command List Generator](https://superfromnd.gitlab.io/Ikemen-cmdlist/), written in Go and released as a standalone executable.

Iguana is a work-in-progress as of this writing (November 2023), and although largely feature-complete is still missing a few features from its web predecessor. [See here](https://github.com/SuperFromND/iguana/issues?q=is%3Aissue+is%3Aopen+label%3Aigoclg-port) for a list of features planned to be ported over.

## Usage
Iguana is a command line tool that requires at least one parameter, `-i`, to specify a command file to load. You can give it either a command (.cmd) file or a definitions (.def) file as input:
```bash
iguana.exe -i path/to/file.cmd
iguana.exe -i path/to/file.def
```

When given a directory, Iguana can also bulk-process all of the .def files it can find within that directory as well as its subfolders:
```bash
iguana.exe -i path/to/folder
```

By default, Iguana will output a file named `movelist.dat` in the same directory as the command file (whether given directly or as part of a .def). If given a .def file, Iguana can also modify it to include support for the movelist.dat file.

Additional options and parameters can be found by running Iguana with either no arguments at all or with the `-h` argument.

## Features

Iguana currently supports the following:
- Motion input and other FG-specialized glyphs
- Power usage annotation
- Customizable header colors
- Fighter Factory-style move labels
- Indirect .cmd processing via .def files
- Bulk processing of entire roster folders
- Automatic .def `movelist` support patching

## Building
Iguana requires at least Go version 1.16 to compile. It only uses a single external module, [go-ini](https://github.com/go-ini/ini), which itself requires Go version 1.13 or later.
```bash
git clone https://github.com/SuperFromND/iguana.git
cd iguana
make
```
The resulting executable will be placed in the `bin` directory.

## Licensing
Iguana's source code, as well as its logo, are [available under the MIT License.](https://raw.githubusercontent.com/SuperFromND/iguana/main/LICENSE) <3