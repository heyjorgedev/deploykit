package main

import (
	"fmt"
	"github.com/heyjorgedev/deploykit"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	dataDir, err := expandPath("./dk-data")
	if err != nil {
		log.Fatal(err)
	}

	app := deploykit.NewWithConfig(deploykit.Config{
		DataDir: dataDir,
		IsDev:   true,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func expandPath(path string) (string, error) {
	// Ignore if path has no leading tilde.
	if path != "~" && !strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return path, nil
	}

	// Fetch the current user to determine the home path.
	u, err := user.Current()
	if err != nil {
		return path, err
	} else if u.HomeDir == "" {
		return path, fmt.Errorf("home directory unset")
	}

	if path == "~" {
		return u.HomeDir, nil
	}
	return filepath.Join(u.HomeDir, strings.TrimPrefix(path, "~"+string(os.PathSeparator))), nil
}

//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
//JORGE000000000000000000000000000000000000LAPA0000000000000
