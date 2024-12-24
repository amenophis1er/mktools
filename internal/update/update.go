// internal/update/update.go
package update

import (
	"encoding/json"
	"fmt"
	"github.com/amenophis1er/mktools/version"
	"net/http"
	"runtime"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func CheckForUpdate() (bool, string, error) {
	current := version.Version
	if current == "dev" {
		return false, "", nil
	}

	resp, err := http.Get("https://api.github.com/repos/amenophis1er/mktools/releases/latest")
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", err
	}

	if release.TagName > current {
		return true, release.TagName, nil
	}

	return false, "", nil
}

func GetUpdateCommand(version string) string {
	switch runtime.GOOS {
	case "darwin", "linux":
		return "brew upgrade mktools"
	case "windows":
		return fmt.Sprintf("curl -L https://github.com/amenophis1er/mktools/releases/download/%s/mktools-windows-amd64.exe -o mktools.exe", version)
	default:
		return fmt.Sprintf("go install github.com/amenophis1er/mktools@%s", version)
	}
}
