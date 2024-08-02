package functions

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SaveFiles(token, courseID string, modules []Module, dirPath string) error {
	courseID = strings.ReplaceAll(courseID, "/", "_")
	var path string
	if dirPath == "" {
		path = filepath.Join(os.Getenv("PWD"), "cursos", courseID)
	} else {
		path = filepath.Join(dirPath, courseID)
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	for i, module := range modules {
		var url string
		if i == 0 {
			url = module.FileURL
		} else {
			url = fmt.Sprintf("%s&token=%s", module.FileURL, token)
		}
		filePath := filepath.Join(path, strings.ReplaceAll(module.FileName, "/", "_"))
		//fmt.Printf("\nDownloading file to: %s\n", filePath)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Couldn't download this file. \n%s\n", err)
			continue
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}
	}

	return nil
}
