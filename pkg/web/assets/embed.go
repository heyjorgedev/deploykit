package assets

import (
	"embed"
	"fmt"
	"github.com/benbjohnson/hashfs"
)

//go:embed js
var embedFS embed.FS

var FS = hashfs.NewFS(embedFS)

func HttpPath(name string) string {
	return fmt.Sprintf("/assets/%s", FS.HashName(name))
}
