package main

import (
	"fmt"
	"log"
	"missingPersons/croatia"
	"missingPersons/dataSource"
)

func main() {
	if err := dataSource.NewDataSource("database", "postgres", "database", "database"); err != nil {
		log.Fatal(err)
	}

	countryMap, err := createCountries([]string{"croatia"})

	if err != nil {
		log.Fatalf("Error occurred while trying to create/find countries: %s. Exiting...", err.Error())
	}

	executions := createExecutions(countryMap)

	for countryName, exec := range executions {
		if err := exec(); err != nil {
			fmt.Printf("Country %s caused an error: %s. Continuing the rest of the countries...", countryName, err.Error())
		}
	}
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
	tx := db.Begin()

	if err := tx.Error; err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, c := range list {
		var country dataSource.Country
		if err := dataSource.FindCountry(c, &country); err != nil {
			country := dataSource.NewCountry("croatia")
			if err := dataSource.SaveCountry(&country); err != nil {
				tx.Rollback()

				return nil, err
			}

			countryList[c] = country

			continue
		}

		countryList[c] = country
	}

	commit := tx.Commit()

	if err := commit.Error; err != nil {
		return nil, err
	}

	return countryList, nil
}
