package cmd

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

// SpotifyKill kills the running Spotify client.
func SpotifyKill() {
	switch runtime.GOOS {
	case "linux":
		if exec.Command("pgrep", "-x", "spotify").Run() == nil {
			exec.Command("pkill", "-x", "spotify").Run()
		}
	case "darwin":
		if exec.Command("sh", "-c", "ps aux | grep 'Spotify' | grep -v grep").Run() == nil {
			exec.Command("pkill", "-x", "Spotify").Run()
		}
	}
}

// SpotifyStart launches the Spotify client.
func SpotifyStart() {
	switch runtime.GOOS {
	case "linux":
		exec.Command(filepath.Join(spotifyPath, "spotify")).Start()
	case "darwin":
		exec.Command("open", "-a", "Spotify", "--args").Start()
	}
}

// SpotifyRestart kills and relaunches the Spotify client.
func SpotifyRestart() {
	SpotifyKill()
	SpotifyStart()
}
