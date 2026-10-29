package cmd

import (
	"path/filepath"

	"github.com/spicetify/cli/src/apply"
	"github.com/spicetify/cli/src/preprocess"
	"github.com/spicetify/cli/src/utils"
)

// Apply blocks ads in Spotify: audio ads, podcast ads and random (mid-roll)
// ads are skipped automatically by injecting Spicetify's player API wrapper and
// the ad skipping extension into the Spotify client.
func Apply(spicetifyVersion string) {
	InitPaths()

	// Copy Spicetify's player API wrapper into the Spotify client.
	if err := utils.CopyFile(
		filepath.Join(utils.GetJsHelperDir(), "spicetifyWrapper.js"),
		filepath.Join(appDestPath, "xpui", "helper")); err != nil {
		utils.PrintError(err.Error())
	}

	preprocess.Start(
		spicetifyVersion,
		spotifyPath,
		appDestPath,
		preprocess.Flag{
			ExposeAPIs: true,
			SpotifyVer: utils.GetSpotifyVersion(prefsPath),
		},
	)

	// Install the ad skipping extension into the Spotify client.
	if extensionPath, err := utils.GetExtensionPath("autoSkipVideo.js"); err == nil {
		if err := utils.CopyFile(extensionPath, filepath.Join(appDestPath, "xpui", "extensions")); err != nil {
			utils.PrintError(err.Error())
		}
	}

	apply.AdditionalOptions(appDestPath, apply.Flag{
		Extension:    []string{"autoSkipVideo.js"},
		SpicetifyVer: spicetifyVersion,
	})
}
