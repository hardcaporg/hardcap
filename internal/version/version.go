package version

import "os"

var (
	// Hostname contains, well, hostname
	Hostname string

	// BuildCommit carries git SHA commit set via -ldflags
	BuildCommit string

	// BuildTime carries date and time in UTC set via -ldflags
	BuildTime string
)

func init() {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown-agent"
	}
	Hostname = h

	if BuildTime == "" {
		BuildTime = "N/A"
	}

	if BuildCommit == "" {
		BuildCommit = "HEAD"
	}
}
