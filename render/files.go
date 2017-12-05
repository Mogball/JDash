package render

import (
	"io/ioutil"
	"log"
	"strings"
)

func GetRenderFiles() []string {
	modelDir := "static/models/"
	files := getFilesInDir(modelDir)
	return files
}

func getFilesInDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Error while reading models directory: %v", err)
		return make([]string, 0)
	}
	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if index := strings.LastIndex(fileName, ".json"); index == len(fileName) - 5 {
				fileNames = append(fileNames, fileName[:index])
			}
		}
	}
	return fileNames
}
