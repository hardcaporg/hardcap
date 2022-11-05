package snip

import "embed"

//go:embed *.ks *.py
var EmbedFS embed.FS
