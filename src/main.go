package main

import (
    "flag"
    "fmt"
    "gopkg.in/ini.v1"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

// basic I/O
var input_file = ""
var output_file = ""

// options
var opt_debug = false
var opt_keep1 = false
var opt_keepai = false

// decorative text for the console
var logo = `
██╗ ██████╗ ██╗   ██╗ █████╗ ███╗   ██╗ █████╗
██║██╔════╝ ██║   ██║██╔══██╗████╗  ██║██╔══██╗
██║██║  ███╗██║   ██║███████║██╔██╗ ██║███████║
██║██║   ██║██║   ██║██╔══██║██║╚██╗██║██╔══██║
██║╚██████╔╝╚██████╔╝██║  ██║██║ ╚████║██║  ██║
╚═╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝
IKEMEN GO Utility for Annotating Nonspecified Attacks
`

var hr = "========================================================"

var version string = "(unknown)"

// stores data from [Command] sections
type Command struct {
    name    string
    command string
    time    int
}

// stores data from [State -1] sections
type Move struct {
    name     string
    triggers []string
}

// combined data from the above two, used when creating movelist data file
type MoveEntry struct {
    name    string
    command string
    held    string
    power   int
}

func check_error(err error) {
    if err != nil {
        panic(err)
    }
}

func prompt() bool {
    // Asks for a yes-or-no prompt, exiting if answered with no
    var confirmation string

    fmt.Printf("[Y/N] ")
    fmt.Scanln(&confirmation)

    if strings.EqualFold(confirmation, "N") {
        os.Exit(0)
        return false
    }

    if strings.EqualFold(confirmation, "Y") {
        return true
    }

    return false
}

func trim_command(input string) string {
    // Strips down Command="cmdinput" to just cmdinput

    if !strings.Contains(strings.ToLower(input), "command") {
        return input
    }

    start := strings.Index(strings.ToLower(input), `command="`) + 9
    end := strings.Index(input[start:], `"`) + start

    return input[start:end]
}

func reverse(input string) string {
    // Flips the given string around

    var output string

    for _, i := range input {
        output = string(i) + output
    }

    return output
}

func remove_element(input []string, index int) []string {
    // Returns an array minus the given element
    j := 0

    for i := range input {
        if i != index {
            input[j] = input[i]
            j++
        }
    }

    output := input[:j]
    return output
}

func merge(input []string) string {
    // Attempt to merge the given strings into one
    // Used in assemble_movelist(), notoriously unstable

    if opt_debug {
        fmt.Println("Attempting to merge commands:", input)
    }

    // determine which one of the commands is the longest
    longest_entry := 0
    var output string

    for e := range input {
        if len(input[e]) > longest_entry {
            longest_entry = e
        }
    }

    // iterate every character in the longest command (to ensure max coverage)
    for c := 0; c < len(input[longest_entry]); c++ {
        var compared_letter string

        for e := range input {
            var current_char string = ""

            // safeguard to prevent OOB read crash
            if len(input[e]) >= len(input[longest_entry]) {
                current_char = input[e][c : c+1]
            }

            // fills the string above with *something* if it's blank
            if compared_letter == "" {
                compared_letter = current_char
                continue
            }

            // if this new character isn't a match with our current one, append it
            if compared_letter != current_char {
                compared_letter += "+" + current_char
            }
        }

        output += compared_letter
    }

    return output
}

func tokenize (output string) string {
    // Tokenizes the command string, replacing each multi-char button input with a single character
    // this is done so that merging them becomes easier

    // ~ indicates releasing a button, unnecessary for us so strip it
    output = strings.ReplaceAll(output, "~", "")

    // > indicates to not press any button betwen previous and next command, unnecessary for us so strip it
    output = strings.ReplaceAll(output, ">", "")

    // $ indicates to read a direction as 4-way, unnecessary for us so strip it
    output = strings.ReplaceAll(output, "$", "")

    // tokenize any held inputs into single characters
    // note the specific characters we're using for tokenizing here are ASCII-compliant
    output = strings.ReplaceAll(output, "/DF", "!")
    output = strings.ReplaceAll(output, "/DB", "@")
    output = strings.ReplaceAll(output, "/UF", "#")
    output = strings.ReplaceAll(output, "/UB", "$")
    output = strings.ReplaceAll(output, "/D", "%")
    output = strings.ReplaceAll(output, "/F", "^")
    output = strings.ReplaceAll(output, "/U", "&")
    output = strings.ReplaceAll(output, "/B", "*")
    output = strings.ReplaceAll(output, "/a", "(")
    output = strings.ReplaceAll(output, "/b", ")")
    output = strings.ReplaceAll(output, "/c", "<")
    output = strings.ReplaceAll(output, "/x", ">")
    output = strings.ReplaceAll(output, "/y", ";")
    output = strings.ReplaceAll(output, "/z", "'")
    output = strings.ReplaceAll(output, "/s", "{")
    output = strings.ReplaceAll(output, "/d", "?")
    output = strings.ReplaceAll(output, "/w", "=")

    // tokenize regular directional inputs
    output = strings.ReplaceAll(output, "DF", "3")
    output = strings.ReplaceAll(output, "DB", "1")
    output = strings.ReplaceAll(output, "UF", "9")
    output = strings.ReplaceAll(output, "UB", "7")
    output = strings.ReplaceAll(output, "D", "2")
    output = strings.ReplaceAll(output, "F", "6")
    output = strings.ReplaceAll(output, "U", "8")
    output = strings.ReplaceAll(output, "B", "4")

    return output
}

func detokenize(output string) string {
    // Converts command from a MoveEntry string into movelist.dat glyphs

    // strip unnecessary commas and spaces
    output = strings.ReplaceAll(output, ",", "")
    output = strings.ReplaceAll(output, " ", "")

    // detokenize held directions
    output = strings.ReplaceAll(output, "!", "~DF")
    output = strings.ReplaceAll(output, "@", "~DB")
    output = strings.ReplaceAll(output, "#", "~UF")
    output = strings.ReplaceAll(output, "$", "~UB")
    output = strings.ReplaceAll(output, "%", "~D")
    output = strings.ReplaceAll(output, "^", "~F")
    output = strings.ReplaceAll(output, "&", "~U")
    output = strings.ReplaceAll(output, "*", "~B")
    output = strings.ReplaceAll(output, "(", "a")
    output = strings.ReplaceAll(output, ")", "b")
    output = strings.ReplaceAll(output, "<", "c")
    output = strings.ReplaceAll(output, ">", "x")
    output = strings.ReplaceAll(output, ";", "y")
    output = strings.ReplaceAll(output, "'", "z")
    output = strings.ReplaceAll(output, "{", "s")
    output = strings.ReplaceAll(output, "?", "d")
    output = strings.ReplaceAll(output, "=", "w")

    // TODO: detokenize option to convert common motion inputs (236 -> _QCF etc etc)

    // double-taps
    output = strings.ReplaceAll(output, "66", "_XFF")
    output = strings.ReplaceAll(output, "44", "_XBB")

    // detokenize regular directions and buttons
    output = strings.ReplaceAll(output, "3", "_DF")
    output = strings.ReplaceAll(output, "1", "_DB")
    output = strings.ReplaceAll(output, "9", "_UF")
    output = strings.ReplaceAll(output, "7", "_UB")
    output = strings.ReplaceAll(output, "2", "_D")
    output = strings.ReplaceAll(output, "6", "_F")
    output = strings.ReplaceAll(output, "8", "_U")
    output = strings.ReplaceAll(output, "4", "_B")

    output = strings.ReplaceAll(output, "a", "^A")
    output = strings.ReplaceAll(output, "b", "^B")
    output = strings.ReplaceAll(output, "c", "^C")
    output = strings.ReplaceAll(output, "x", "^X")
    output = strings.ReplaceAll(output, "y", "^Y")
    output = strings.ReplaceAll(output, "z", "^Z")
    output = strings.ReplaceAll(output, "s", "^S")
    output = strings.ReplaceAll(output, "d", "^D")
    output = strings.ReplaceAll(output, "w", "^W")

    // other special glyphs
    output = strings.ReplaceAll(output, "+", "_+")

    return output
}


func detect_ai_command(c Command) bool {
    // Checks if the given input command is a WinMUGEN-style AI command
    // WinMUGEN did not yet include a trigger for checking if a character was AI-controlled,
    // so a workaround many authors used was to make commands impossible for a human to input,
    // the idea being that MUGEN's AI could just "activate" it on a whim without actually inputting the command

    // tokenizes the input so 1 input = 1 character
    // we also strip out any formatting/padding
    cmd_str := tokenize(c.command)
    cmd_str = strings.ReplaceAll(cmd_str, ",", "")
    cmd_str = strings.ReplaceAll(cmd_str, " ", "")
    cmd_str = strings.ReplaceAll(cmd_str, "+", "")

    // the length of this string should be roughly equal to the number of inputs the game has to detect for this command
    button_count := len(cmd_str)

    // a zero-time command is physically impossible to do
    if (c.time == 0) {return true}

    // checks if there's more buttons to press than the number of frames where inputs are read as part of a command
    if (c.time < button_count) {return true}

    return false
}

func scrape_commands(input *ini.File) []Command {
    // Returns array of command-structs created from the given INI
    // This should *only* parse sections named "Command" (case insensitive)

    if opt_debug {
        fmt.Println("Scraping command entries...", "\n"+hr)
    }
    var cmdlist []Command

    for s := range input.Sections() {
        var sect_name = input.Sections()[s].Name()

        if strings.EqualFold(sect_name, "Command") {
            var cmd Command

            // fallback/default value
            // TODO: get actual value from Defaults section
            cmd.time = 15

            for k := range input.Sections()[s].Keys() {
                var key_name = input.Sections()[s].KeyStrings()[k]

                if strings.EqualFold(key_name, "name") {
                    cmd.name = input.Sections()[s].Key(key_name).String()
                }

                if strings.EqualFold(key_name, "command") {
                    // strip out all spaces from the command internally, done to make comparisons easier (see assemble_move_table() below)
                    cmd.command = strings.ReplaceAll(input.Sections()[s].Key(key_name).String(), " ", "")
                }

                // amount of time the command must be performed in, used to check for possible AI commands
                if strings.EqualFold(key_name, "time") {
                    cmd.time, _ = input.Sections()[s].Key(key_name).Int()
                }
            }

            if opt_debug {
                fmt.Println("Found command:", cmd)
            }
            cmdlist = append(cmdlist, cmd)
        }

        // Statedef -1 marks the end of [Command] blocks, so once we reach that, end the loop
        if strings.EqualFold(sect_name, "Statedef -1") {
            if opt_debug {
                fmt.Println("Reached Statedef -1. Ending command scraping...")
            }
            break
        }
    }

    return cmdlist
}

func scrape_moves(input *ini.File) []Move {
    // Returns array of move-structs created from the given INI
    // This should *only* parse sections after the [Statedef -1] section

    if opt_debug {
        fmt.Println("Scraping move state controllers...", "\n"+hr)
    }
    var moves []Move
    var statedef_reached = false

    for s := range input.Sections() {
        var sect_name = input.Sections()[s].Name()

        if statedef_reached {
            var move Move
            var is_changestate bool = true
            var has_command = false

            // trims "State -1," and then trims any whitespace
            move.name = strings.TrimSpace(sect_name[strings.Index(sect_name, ",")+1:])

            // scans over every key in the state controller
            for k := range input.Sections()[s].Keys() {
                var key_name = input.Sections()[s].KeyStrings()[k]

                // this extra range loop is necessary to parse "shadow keys", what the INI lib calls multiple keys with identical names (e.g. multiple triggeralls)
                for v := range input.Sections()[s].Key(key_name).StringsWithShadows("\n") {
                    var key_value = input.Sections()[s].Key(key_name).StringsWithShadows("\n")[v]

                    // checks the type field to make sure that this actually is a move
                    // if it isn't, break the loop and move onto the next section
                    if strings.EqualFold(key_name, "type") {
                        if !strings.EqualFold(key_value, "ChangeState") {
                            if opt_debug {
                                fmt.Println("Move", move.name, "detected as a non-ChangeState type. Discarding...")
                            }
                            is_changestate = false
                            break
                        }
                    }

                    // checks the keys to see if they contain a command or a trigger related to power
                    // the latter is what determines if a move is a hyper
                    var trigger_has_command bool = strings.Contains(strings.ToLower(key_value), strings.ToLower("command"))
                    var trigger_has_power bool = strings.Contains(strings.ToLower(key_value), strings.ToLower("power"))

                    // check to make sure the command doesn't have the "not equals" operator (effectively rendering it null for us)
                    if trigger_has_command {
                        if strings.Contains(strings.ReplaceAll(strings.ToLower(key_value), " ", ""), strings.ToLower("!=")) {
                            if opt_debug {
                                fmt.Println("Move trigger", key_value, "detected as not-equals op. Discarding...")
                            }
                            trigger_has_command = false
                        }

                        has_command = true
                    }

                    if trigger_has_command || trigger_has_power {
                        // we also strip out spaces for the triggers here, to make command detection easier
                        move.triggers = append(move.triggers, strings.ReplaceAll(key_value, " ", ""))
                    }
                }
            }

            if is_changestate && has_command {
                if opt_debug {
                    fmt.Println("Found move:", move)
                }
                moves = append(moves, move)
            }
        }

        // Statedef -1 is where moves start being defined, so we ignore everything before that point as an optimization
        if strings.EqualFold(sect_name, "Statedef -1") {
            statedef_reached = true
        }
    }

    return moves
}

func assemble_move_table(commands []Command, moves []Move) []MoveEntry {
    // Takes commands and moves, then converts them to entries in an array
    // The returned array is then formatted by format_move_table() and saved to disk

    if opt_debug {
        fmt.Println("Assembling move table...", "\n"+hr)
    }
    var movelist []MoveEntry

    for m := range moves {
        var mv MoveEntry
        var cmds []string
        var cmd_str string

        mv.name = moves[m].name
        if opt_debug {
            fmt.Println("Reading move:", mv.name)
        }

        for t := range moves[m].triggers {
            move_command := trim_command(moves[m].triggers[t])

            for c := range commands {
                // checks if the current trigger has a corresponding command
                if move_command == commands[c].name {
                    if (!opt_keepai && detect_ai_command(commands[c])) {
                        if opt_debug {
                            fmt.Println("Command detected as AI-only:", move_command)
                        }
                    } else {
                        if opt_debug {
                            fmt.Println("Tokenizing string:", commands[c].command)
                        }

                        command_text := tokenize(commands[c].command)

                        if opt_debug {
                            fmt.Println("Tokenized:", command_text)
                        }
                        cmds = append(cmds, command_text)
                    }
                }
            }

            // detect power requirements and pack them into the move entry
            if strings.Contains(strings.ToLower(move_command), "power") {
                if opt_debug {
                    fmt.Println("Power requirement detected:", move_command)
                }

                trim_amount := 6
                if strings.Contains(strings.ToLower(move_command), ">=") {
                    trim_amount++
                }

                pow, err := strconv.ParseInt(move_command[trim_amount:], 10, 16)
                if err != nil {
                    pow = 0
                }

                mv.power = int(pow)
            }
        }

        // separate the held-input tokens from the regular tokens
        // held inputs are traditionally kept at the beginning of a command, so we put them into an array to read off later
        var holding []string
        var nonholding []string
        token_regex, _ := regexp.Compile("[!@#$%^&*()[\\];'.~>]")

        for e := range cmds {
            if len(cmds[e]) == 1 && token_regex.MatchString(cmds[e][0:1]) {
                holding = append(holding, cmds[e])
            } else {
                nonholding = append(nonholding, cmds[e])
            }
        }
        if len(nonholding) == 0 {
            // do nothing
        } else if len(nonholding) == 1 {
            // there's only one command, so just use it verbatim
            cmd_str = nonholding[0]
        } else {
            // heuristic: if there's only two non-held commands that are identical flipped, return one of them
            // this specifically is meant to catch KFM's blocking inputs among other similar commands
            if len(nonholding) == 2 && nonholding[0] == reverse(nonholding[1]) {
                if opt_debug {
                    fmt.Println("Commands detected to be identical when mirrored. Returning one of them...")
                }
                cmd_str = nonholding[0]
            } else {
                // attempt to merge the commands (probably the most inaccurate part of the entire program, thus avoided when possible)
                cmd_str = merge(nonholding)
            }
        }

        if holding != nil {
            if opt_debug {
                fmt.Println("Held inputs found:", holding)
            }
            mv.held = fmt.Sprintf("%s", holding)
        }

        // fallback if, for whatever reason, we end up with no inputs at all
        if cmd_str == "" && holding == nil {
            if opt_debug {
                fmt.Println("Move", mv.name, "has no usable commands. Discarding...")
            }
            continue
        }

        mv.command = cmd_str
        movelist = append(movelist, mv)
    }

    return movelist
}

func format_move_table(move_table []MoveEntry) string {
    // Takes a given array of moves and formats it into movelist.dat's formatting
    // What this function returns is then saved to disk

    if opt_debug {
        fmt.Println("Formatting move table...", "\n"+hr)
    }

    special_list := "<#" + "f0f000" + ">:Special Moves:</>\n"
    hypers_list := "<#" + "f0f000" + ">:Hyper Moves:</>\n"

    for i := range move_table {
        var entry string

        // remove one-button, non-hyper commands
        if !opt_keep1 && len(move_table[i].command) == 1 && move_table[i].power == 0 {
            continue
        }

        cmd := detokenize(move_table[i].held) + detokenize(move_table[i].command)

        entry = move_table[i].name

        if move_table[i].power != 0 {
            entry += " <#" + "bebebe" + ">(" + strconv.Itoa(move_table[i].power) + ")</>"
            hypers_list += entry + "\t\t\t" + cmd + "\n"
            continue
        }

        special_list += entry + "\t\t\t" + cmd + "\n"
    }

    return special_list + "\n" + hypers_list
}

func Convert(path string) string {
    // wrapper function that does all of the actual work
    // takes a file path as input and returns a movelist.dat as a string
    // note the capital, we actually export this function for use by other programs (wink wink)

    // loads the file, then parses it
    if opt_debug {
        fmt.Println("Reading input file...")
    }
    file_data, err := os.ReadFile(path)
    check_error(err)

    if opt_debug {
        fmt.Println("Parsing as INI data...")
    }
    parsed_ini, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true, AllowShadows: true, SkipUnrecognizableLines: true}, file_data)
    check_error(err)

    // parse sections into dedicated structs
    commands := scrape_commands(parsed_ini)
    moves := scrape_moves(parsed_ini)

    // combine the parsed data into a list of move names and command inputs
    move_table := assemble_move_table(commands, moves)

    // format the movelist we just made into the movelist.dat format
    // (see https://github.com/ikemen-engine/Ikemen-GO/wiki/Miscellaneous-Info#movelists)
    return format_move_table(move_table)
}

func main() {
    // Re-define output of -h
    flag.Usage = func() {
        footer := `Program written by SuperFromND.
Distributed under the MIT license.
`
        // Print cool logo and separation bars
        fmt.Printf(logo + "version " + version + "\n" + footer + hr + "\n")
        fmt.Printf("\nCommand arguments for IGUANA:\n")
        flag.PrintDefaults()
        fmt.Printf("\n")
        os.Exit(0)
    }

    flag.StringVar(&input_file, "i", "", "command file to parse (required)")
    flag.StringVar(&output_file, "o", "movelist.dat", "output filename, excluding path")
    flag.BoolVar(&opt_debug, "d", false, "enables debug logging")
    flag.BoolVar(&opt_keep1, "keep1", false, "preserve one-button, non-hyper moves")
    flag.BoolVar(&opt_keepai, "keepai", false, "preserve move commands detected as AI-only")

    flag.Parse()

    // check if any arguments are present
    if len(os.Args[1:]) == 0 {
        flag.Usage()
        os.Exit(0)
    }

    // check if an input file has been given
    if len(input_file) == 0 {
        fmt.Printf("No input files given. Syntax is 'iguana -i command.cmd'\n")
        os.Exit(0)
    }

    info, err := os.Stat(input_file)
    check_error(err)

    if info.IsDir() {
        // input is a directory, so this means we're in batch mode
        prompt_msg := `
Iguana has been given a directory as input and is in batch mode.
It will attempt to process every command file in every sub-folder in this directory.
Making a backup of this folder is recommended before continuing.
Are you sure you want to continue? `

        fmt.Printf(prompt_msg)

        if prompt() {
            var cmd_file_list []string

            filepath.Walk(input_file, func(path string, info os.FileInfo, err error) error {
                if filepath.Ext(path) == ".cmd" {cmd_file_list = append(cmd_file_list, path)}
                return nil
            })

            for i := range cmd_file_list {
                fmt.Println("Converting file: " + cmd_file_list[i])
                movelist := Convert(cmd_file_list[i])

                if opt_debug {
                    fmt.Println("Dump of movelist:\n" + movelist)
                } else {
                    path := filepath.Dir(cmd_file_list[i]) + "/" + output_file
                    err := os.WriteFile(path, []byte(movelist), 0666)
                    check_error(err)
                }
            }
        }

    } else {
        // check to make sure the input file has the right extension
        if filepath.Ext(input_file) != ".cmd" {
            fmt.Printf("Input file is not a command (.cmd) file.\n")
            os.Exit(0)
        }

        // make a note if debug logging is on
        if opt_debug {
            fmt.Println("Debug logging enabled.")
        }

        // at this point, we know we have a file, so try to do stuff with it
        movelist := Convert(input_file)

        if opt_debug {
            fmt.Println("Dump of movelist:\n" + movelist)
        } else {
            path := filepath.Dir(input_file) + "/" + output_file
            fmt.Println("Saving to path: " + path)
            err := os.WriteFile(path, []byte(movelist), 0666)
            check_error(err)
        }
    }
}
