package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spicetify/cli/src/apply"
	"github.com/spicetify/cli/src/preprocess"
	"github.com/spicetify/cli/src/utils"
)

const lastAppliedVersionFileName = "last-applied-version"

// Apply blocks ads in Spotify: audio ads, podcast ads and random (mid-roll)
// ads are skipped automatically by injecting Spicetify's player API wrapper and
// the ad skipping extension into the Spotify client.
// It re-applies the ad blocking automatically when Spotify got updated.
// The caller is responsible for restarting the Spotify client.
func Apply(spicetifyVersion string) {
	ApplyIfRequired(spicetifyVersion)
}

// ApplyIfRequired applies the ad blocking when it is not applied yet or when
// Spotify got updated since the last apply, so users don't have to redo it.
func ApplyIfRequired(spicetifyVersion string) {
	InitPaths()

	spotifyVersion := utils.GetSpotifyVersion(prefsPath)
	if spotifyVersion == getLastAppliedSpotifyVersion() {
		return
	}

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
			SpotifyVer: spotifyVersion,
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

	setLastAppliedSpotifyVersion(spotifyVersion)
}

// getLastAppliedSpotifyVersion reads the Spotify version the ad blocking was
// applied with. Empty string when the ad blocking was never applied.
func getLastAppliedSpotifyVersion() string {
	versionPath := filepath.Join(utils.GetStateFolder("AdBlock"), lastAppliedVersionFileName)
	content, err := os.ReadFile(versionPath)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(content))
}

// setLastAppliedSpotifyVersion stores the Spotify version the ad blocking was
// applied with.
func setLastAppliedSpotifyVersion(version string) {
	versionPath := filepath.Join(utils.GetStateFolder("AdBlock"), lastAppliedVersionFileName)
	if err := os.WriteFile(versionPath, []byte(version), 0700); err != nil {
		utils.PrintError(err.Error())
	}
}
