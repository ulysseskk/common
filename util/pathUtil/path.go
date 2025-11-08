package pathUtil

import (
	"os"
	"path"
	"strings"
)

func TransHomePrefixPathToAbsolutePath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", homePath, 1)
	}
	return path, nil
}

func GetAbsolutePathForCurrentProcess(p string) (string, error) {
	if path.IsAbs(p) {
		return p, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(wd, p), nil
}
