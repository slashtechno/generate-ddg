// From the guide: "In this code, the //go:embed "html" "static" directive embeds the contents of the assets/html and assets/static directories into the files variable, which is an embed.FS rooted in the assets directory."

package assets

import (
	"embed"
	"io/fs"
)

// https://pkg.go.dev/embed#hdr-File_Systems
//
//go:embed "html" "static"
var files embed.FS

var (
	HTMLFiles   = sub(files, "html")
	StaticFiles = sub(files, "static")
)

func sub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
