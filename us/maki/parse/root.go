package parse

import (
	"path/filepath"

	"github.com/darabuchi/utils"
)

var rootDir string

func Root() string {
	if rootDir != "" {
		return rootDir
	}

	current := utils.GetPwd()

	for !utils.IsDir(filepath.Join(current, "proto")) || current == "" {
		current = filepath.Dir(current)
	}

	rootDir = current
	return rootDir
}

func ProtoDir() string {
	return filepath.Join(Root(), "proto")
}
