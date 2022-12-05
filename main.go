package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"missingPersons/dataSource"
	"missingPersons/download"
	croatia2 "missingPersons/internal/croatia"
	"missingPersons/logger"
	"os"
	"sync"
)

func main() {
	LoadEnv()
	fmt.Println("Creating loggers...")
	if err := logger.BuildLoggers([]string{"croatia", "serbia"}); err != nil {
		log.Fatalf("Unable to build loggers: %s\n", err.Error())
	}
	fmt.Println("Loggers created!")

	fmt.Println("Creating image directory...")
	createImageDirIfNotExists()
	fmt.Println("Image directory created!")

	if err := dataSource.NewDataSource("database", "postgres", "database", "database"); err != nil {
		log.Fatalf("Cannot connect to postgres database: %s", err.Error())
	}

	fmt.Println("Creating countries if they do not exist...")
	countryMap, err := createCountries([]string{"croatia"})
	if err != nil {
		log.Fatalf("Error occurred while trying to create/find countries: %s. Exiting...", err.Error())
	}

	fmt.Println("Countries created or fetched. Continuing...\n")

	run(countryMap)

	fmt.Println("")
	fmt.Println("Process finished!")
}

func run(countryMap map[string]dataSource.Country) {
	executions := createExecutions(countryMap)

	wg := &sync.WaitGroup{}
	for countryName, exec := range executions {
		wg.Add(1)
		go func(exec func() error, wg *sync.WaitGroup, countryName string) {
			if err := exec(); err != nil {
				fmt.Printf("Country %s caused an error: %s. Continuing the rest of the countries...", countryName, err.Error())
			}

			wg.Done()
		}(exec, wg, countryName)
	}

	wg.Wait()
}

func createExecutions(countryMap map[string]dataSource.Country) map[string]func() error {
	list := make(map[string]func() error, 0)

	list["Croatia"] = func() error {
		people, err := croatia2.StartScrapping()

		if err != nil {
			return err
		}

		fmt.Printf("Croatia: Found %d people\n", len(people))
		info, err := croatia2.SaveCountry(people, countryMap["croatia"], download.NewFsImageSaver())

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

	return list
}

func createCountries(list []string) (map[string]dataSource.Country, error) {
	countryList := make(map[string]dataSource.Country, 0)

	db := dataSource.DB()

	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, c := range list {
			var country dataSource.Country
			if err := db.Where("name = ?", c).First(&country).Error; err != nil {
				country := dataSource.NewCountry("croatia")

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

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}
