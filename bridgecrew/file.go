package bridgecrew

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// loadFileContent returns contents of a file in a given path
func loadFileContent(v string) ([]byte, error) {
	//nolint
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
