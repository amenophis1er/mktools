// cmd/version.go
package cmd

import (
	"fmt"
	"runtime"

	"github.com/amenophis1er/mktools/internal/update"
	"github.com/amenophis1er/mktools/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Show mktools version information",
        Run: func(cmd *cobra.Command, args []string) {
            // Print version info
            fmt.Printf("mktools %s (%s)\n", version.Version, version.Commit)
            fmt.Printf("Built with %s on %s\n", runtime.Version(), version.Date)

            // Note for development builds
            if version.Version == "main" {
                fmt.Println("\nNote: This is a development build. For release versions, install via:")
                fmt.Println("brew tap amenophis1er/homebrew-mktools && brew install mktools")
                fmt.Println("or download from https://github.com/amenophis1er/mktools/releases")
                return
            }

            // Check for updates (only for release versions)
            hasUpdate, newVersion, err := update.CheckForUpdate()
            if err != nil {
                return // Silently fail update check
            }

            if hasUpdate {
                fmt.Printf("\nNew version available: %s\n", newVersion)
                fmt.Printf("To update, run: %s\n", update.GetUpdateCommand(newVersion))
            }
        },
    }
}
