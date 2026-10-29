package apply

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spicetify/cli/src/utils"
)

// Flag enables/disables additional feature
type Flag struct {
	// Extension is the list of extension files to be injected into Spotify.
	Extension []string
	// SpicetifyVer is the Spicetify version to be recorded in the client's config.
	SpicetifyVer string
}

// AdditionalOptions injects extensions into Spotify's index.html.
func AdditionalOptions(appsFolderPath string, flags Flag) {
	htmlPath := filepath.Join(appsFolderPath, "xpui", "index.html")

	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		return
	}

	var extensionsHTML strings.Builder
	extensionsHTML.WriteByte('\n')

	if flags.SpicetifyVer != "" {
		var extList strings.Builder
		for _, ext := range flags.Extension {
			fmt.Fprintf(&extList, `"%s",`, ext)
		}

		fmt.Fprintf(&extensionsHTML, `<script>
			Spicetify.Config={};
			Spicetify.Config["version"]="%s";
			Spicetify.Config["extensions"] = [%s];
		</script>
		`, flags.SpicetifyVer, extList.String())
	}

	for _, v := range flags.Extension {
		if strings.HasSuffix(v, ".mjs") {
			fmt.Fprintf(&extensionsHTML, "<script defer type='module' src='extensions/%s'></script>\n", v)
		} else {
			fmt.Fprintf(&extensionsHTML, "<script defer src='extensions/%s'></script>\n", v)
		}
	}

	utils.ModifyFile(htmlPath, func(content string) string {
		utils.Replace(
			&content,
			`<script defer="defer" src="/xpui-snapshot\.js"></script>`,
			func(submatches ...string) string {
				return `<script defer="defer" src="/xpui-modules.js"></script><script defer="defer" src="/xpui-snapshot.js"></script>`
			})

		utils.Replace(
			&content,
			`</body>`,
			func(submatches ...string) string {
				return fmt.Sprintf("%s%s", extensionsHTML.String(), submatches[0])
			})

		return content
	})
}
