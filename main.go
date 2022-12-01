package main

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"missingPersons/croatia"
	"missingPersons/dataSource"
	"sync"
)

func main() {
	if err := dataSource.NewDataSource("database", "postgres", "database", "database"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating countries if they do not exist...")
	countryMap, err := createCountries([]string{"croatia"})
	fmt.Println("Countries created or fetched. Continuing...\n")

	if err != nil {
		log.Fatalf("Error occurred while trying to create/find countries: %s. Exiting...", err.Error())
	}

	executions := createExecutions(countryMap)

	wg := &sync.WaitGroup{}
	for countryName, exec := range executions {
		wg.Add(1)
		go func(exec func() error, wg *sync.WaitGroup) {
			if err := exec(); err != nil {
				fmt.Printf("Country %s caused an error: %s. Continuing the rest of the countries...", countryName, err.Error())
			}

			wg.Done()
		}(exec, wg)
	}

	wg.Wait()

	fmt.Println("")
	fmt.Println("Process finished!")
}

func createExecutions(countryMap map[string]dataSource.Country) map[string]func() error {
	list := make(map[string]func() error, 0)

	list["Croatia"] = func() error {
		pages := croatia.StartScrapping()
		fmt.Printf("Croatia: Found %d pages\n", len(pages))
		info, err := croatia.SaveCountry(pages, countryMap["croatia"])

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
			if err := dataSource.FindCountry(c, &country); err != nil {
				country := dataSource.NewCountry("croatia")
				if err := dataSource.SaveCountry(&country); err != nil {
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
