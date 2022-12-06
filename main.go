package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"missingPersons/dataSource"
	"missingPersons/download"
	"missingPersons/internal/croatia"
	"missingPersons/internal/serbia"
	"os"
)

func main() {
	loadEnv()

	runDb()
	runLoggers()
	runCreateImageDir()
	runScrappers(runCountries())

	fmt.Println("")
	fmt.Println("Process finished!")
}

func createExecutions(countryMap map[string]dataSource.Country) map[string]func() error {
	list := make(map[string]func() error, 0)

	list["Croatia"] = func() error {
		people, err := croatia.StartScrapping()

		if err != nil {
			return err
		}

		fmt.Printf("Croatia: Found %d people\n", len(people))
		info, err := croatia.SaveCountry(people, countryMap["croatia"], download.NewFsImageSaver())

		if err != nil {
			return err
		}

		fmt.Printf(`
Croatia:
    Created entries: %d
    Updated entries: %d
    Deleted entries: %d
`, info.CreatedCount, info.UpdatedCount, info.DeletedCount)

		return nil
	}

	list["Serbia"] = func() error {
		people, err := serbia.StartScrapping()

		if err != nil {
			return err
		}

		fmt.Printf("Serbia: Found %d people\n", len(people))
		info, err := serbia.SaveCountry(people, countryMap["serbia"])

		if err != nil {
			return err
		}

		fmt.Printf(`
Serbia:
    Created entries: %d
    Updated entries: %d
    Deleted entries: %d
`, info.CreatedCount, info.UpdatedCount, info.DeletedCount)

		return nil
	}

	return list
}

func createCountries(list []string) (map[string]dataSource.Country, error) {
	countryList := make(map[string]dataSource.Country, 0)

	db := dataSource.DB()

	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, c := range list {
			var country dataSource.Country
			if err := db.Where("name = ?", c).First(&country).Error; err != nil {
				country := dataSource.NewCountry(c)

				if err := db.Create(&country).Error; err != nil {
					return err
				}

				countryList[c] = country

				continue
			}

			countryList[c] = country
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return countryList, nil
}

func createImageDirIfNotExists() {
	dir := os.Getenv("IMAGE_DIRECTORY")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fsErr := os.MkdirAll(dir, os.ModePerm)

		if fsErr != nil {
			log.Fatal(fmt.Sprintf("Cannot create images directory: %s", fsErr.Error()))
		}
	}
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}
