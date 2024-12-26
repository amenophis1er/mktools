// version/version.go
package version

var (
    // Overridden during build via -ldflags
    Version = "main"
    Commit  = "none"
    Date    = "unknown"
)