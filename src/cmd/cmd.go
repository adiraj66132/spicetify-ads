package cmd

import (
	"os"
	"path/filepath"

	"github.com/spicetify/cli/src/utils"
)

var (
	spotifyPath string
	prefsPath   string
	appDestPath string
)

// InitPaths tries to detect Spotify client's installation paths.
// Spotify is expected on Linux or macOS only.
func InitPaths() {
	spotifyPath = utils.FindAppPath()

	if len(spotifyPath) == 0 {
		utils.PrintError(`Cannot detect Spotify location. Please install Spotify or point the CLI to it.`)
		os.Exit(1)
	}

	prefsPath = utils.FindPrefFilePath()

	if len(prefsPath) == 0 {
		utils.PrintError(`Cannot detect Spotify "prefs" file location.`)
		os.Exit(1)
	}

	appDestPath = filepath.Join(spotifyPath, "Apps")
}

// GetSpotifyPath returns location of Spotify client
func GetSpotifyPath() string {
	return spotifyPath
}
