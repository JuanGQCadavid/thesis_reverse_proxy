package utils

import (
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
)

func FromFilePathToStruct(path string, data interface{}) error {
	fileOnBytes, err := os.ReadFile(path)

	if err != nil {
		return errors.Join(fmt.Errorf("err while reading file located at ", path), err)
	}

	if err := yaml.Unmarshal(fileOnBytes, data); err != nil {
		return errors.Join(fmt.Errorf("err while unmarshalling file"), err)
	}
	return nil
}
