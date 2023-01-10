//adapted from 'micro's way of accomplishing this task https://github.com/zyedidia/micro

package config

import (
	"embed"
	"path/filepath"
	"strings"
)

//go:embed colorschemes help syntax plugins
var runtime embed.FS

func fixPath(name string) string {
	return strings.TrimLeft(filepath.ToSlash(name), "runtime/")
}

// AssetDir lists file names in folder
func AssetDir(name string) ([]string, error) {
	name = fixPath(name)
	entries, err := runtime.ReadDir(name)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries), len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names, nil
}

// Asset returns a file content
func Asset(name string) ([]byte, error) {
	name = fixPath(name)
	return runtime.ReadFile(name)
}
