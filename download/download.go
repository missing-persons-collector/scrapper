package download

import (
	"errors"
	"fmt"
	"io"
	"missingPersons/httpClient"
	"net/http"
	"os"
	"strings"
)

type ImageSaver interface {
	Save(URL string, id string) error
}

type fsImageSaver struct{}

func (i fsImageSaver) Save(URL string, id string) error {
	fileName := CreateImageName(URL, id)

	_, err := downloadAndSaveImage(URL, fileName)

	return err
}

func NewFsImageSaver() ImageSaver {
	return fsImageSaver{}
}

func downloadAndSaveImage(URL, fileName string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath := fmt.Sprintf("%s/%s/%s", path, os.Getenv("IMAGE_DIRECTORY"), fileName)

	response, err := httpClient.SendRequest(URL)

	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Request returned non 200 for %s", URL))
	}

	file, err := os.Create(absPath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func CreateImageName(url string, id string) string {
	split := strings.Split(url, ".")
	return fmt.Sprintf("%s.%s", id, split[len(split)-1])
}
