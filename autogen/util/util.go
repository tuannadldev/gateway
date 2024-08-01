package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetAllProtoFiles(path string) []string {
	var protoPaths []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".proto") {
			protoPaths = append(protoPaths, path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("failed to get all proto files, return whatever found, error:", err)
	}

	return protoPaths
}

func ReadFileAsString(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CreateFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	// folder has already existed
	return nil
}

func LowerFirstChar(input string) string {
	if len(input) <= 1 {
		return input
	}

	return strings.ToLower(input[:1]) + input[1:]
}

func UpperFirstChar(input string) string {
	if len(input) <= 1 {
		return input
	}

	return strings.ToUpper(input[:1]) + input[1:]
}

func ConvertSliceStrToUCWord(input []string) string {
	mapString := ""
	for _, str := range input {
		mapString += UpperFirstChar(str)
	}
	return mapString
}
