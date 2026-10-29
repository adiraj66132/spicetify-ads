package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/spicetify/cli/src/cmd"
)

var version string

var (
	quiet     = false
	noRestart = false
	commands  []string
)

func init() {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		log.Println("spicetify only supports Linux and macOS.")
		os.Exit(1)
	}

	if version == "" {
		version = "Dev"
	}

	log.SetFlags(0)
}

func main() {
	parseFlags := true
	for _, v := range os.Args[1:] {
		if parseFlags && v == "--" {
			parseFlags = false
			continue
		}

		if parseFlags && len(v) > 0 && v[0] == '-' {
			switch v {
			case "-q", "--quiet":
				quiet = true
			case "-h", "--help":
				help()
				os.Exit(0)
			case "-v", "--version":
				log.Println(version)
				os.Exit(0)
			case "-n", "--no-restart":
				noRestart = true
			default:
				log.Println("Unknown flag \"" + v + "\"")
				help()
				os.Exit(1)
			}
			continue
		}

		commands = append(commands, v)
	}

	if quiet {
		log.SetOutput(io.Discard)
		pterm.DisableOutput()
	}

	if len(commands) < 1 {
		help()
		os.Exit(0)
	}

	switch commands[0] {
	case "apply":
		cmd.Apply(version)
		if !noRestart {
			cmd.SpotifyRestart()
		}
	case "restart":
		cmd.SpotifyRestart()
	default:
		log.Println("Command \"" + commands[0] + "\" not found.")
		log.Println(`Run "spicetify -h" for a list of valid commands.`)
		os.Exit(1)
	}
}

func help() {
	log.Println("spicetify v" + version)
	log.Println("USAGE")
	log.Println("  spicetify apply    Block ads in Spotify (audio ads, podcast ads and random ads)")
	log.Println("  spicetify restart  Restart the Spotify client")
	log.Println()
	log.Println("FLAGS")
	log.Println("  -q, --quiet         Quiet mode (no output).")
	log.Println("  -n, --no-restart    Do not restart Spotify after running 'spicetify apply'.")
	log.Println("  -h, --help           Print this help text and quit.")
	log.Println("  -v, --version        Print version number and quit.")
}
