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

func DownloadAndSaveImage(URL, fileName string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath := fmt.Sprintf("%s/%s/%s", path, "images", fileName)

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

func CreateAbsPath(fileName string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath := fmt.Sprintf("%s/%s/%s", path, "images", fileName)

	return absPath, nil
}

func ImageExists(fileName string) (bool, error) {
	absPath, err := CreateAbsPath(fileName)

	if err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)

	if os.IsNotExist(err) {
		return false, nil
	}

	return os.IsExist(err), nil
}
