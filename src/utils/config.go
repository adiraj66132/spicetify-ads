package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// FindAppPath finds Spotify location in various possible places
// of each platform and returns it.
// Returns blank string if none of default locations exists.
func FindAppPath() string {
	switch runtime.GOOS {
	case "linux":
		return linuxApp()
	case "darwin":
		return darwinApp()
	}

	return ""
}

// FindPrefFilePath finds Spotify "prefs" file location
// in various possible places of each platform and returns it.
// Returns blank string if none of default locations exists.
func FindPrefFilePath() string {
	switch runtime.GOOS {
	case "linux":
		return linuxPrefs()
	case "darwin":
		return darwinPrefs()
	}

	return ""
}

func linuxApp() string {
	path, err := exec.Command("whereis", "-b", "spotify").Output()

	if err == nil {
		pathString := strings.Replace(string(path), "spotify: ", "", 1)
		pathString = strings.Replace(pathString, "\n", "", -1)
		binList := strings.Split(pathString, " ")

		for _, v := range binList {
			bin := v

			stat, err := os.Lstat(bin)
			if err != nil {
				continue
			}

			if (stat.Mode() & os.ModeSymlink) != 0 {
				binDest, err := os.Readlink(v)

				if err != nil {
					continue
				}

				bin = binDest
			}

			bin = filepath.Dir(bin)

			if _, err := os.Stat(filepath.Join(bin, "Apps")); err == nil {
				return bin
			}
		}
	}

	potentialList := []string{
		"/opt/spotify/",
		"/opt/spotify/spotify-client/",
		"/usr/share/spotify/",
		"/usr/libexec/spotify/",
		"/usr/share/spotify-client/",
		"/var/lib/flatpak/app/com.spotify.Client/x86_64/stable/active/files/extra/share/spotify/",
		"$HOME/.local/share/flatpak/app/com.spotify.Client/x86_64/stable/active/files/extra/share/spotify/",
		"$HOME/.local/share/spotify-launcher/install/usr/share/spotify/",
	}

	for _, v := range potentialList {
		_, err := os.Stat(filepath.Join(ReplaceEnvVarsInString(v), "Apps"))
		_, err2 := os.Stat(filepath.Join(ReplaceEnvVarsInString(v), "spotify"))
		if err == nil && err2 == nil {
			return v
		}
	}

	return ""
}

func linuxPrefs() string {
	dotConfig := os.Getenv("XDG_CONFIG_HOME")

	if len(dotConfig) == 0 {
		dotConfig = filepath.Join(os.Getenv("HOME"), ".config")
	}

	pref := filepath.Join(dotConfig, "spotify", "prefs")
	if _, err := os.Stat(pref); err == nil {
		return pref
	}

	return ""
}

func darwinApp() string {
	path := filepath.Join("/Applications", "Spotify.app", "Contents", "Resources")
	if _, err := os.Stat(path); err == nil {
		return path
	}

	return ""
}

func darwinPrefs() string {
	pref := filepath.Join(os.Getenv("HOME"), "Library/Application Support/Spotify/prefs")
	if _, err := os.Stat(pref); err == nil {
		return pref
	}

	return ""
}
