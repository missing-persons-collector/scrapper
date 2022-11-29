package main

import (
	"log"
	"missingPersons/contract"
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

	saveCountry(croatia.Start(), countryMap["croatia"])
}

func saveCountry(pages []contract.CollectedPage, country dataSource.Country) {
	tx := dataSource.Transaction()
	for _, page := range pages {
		data := page.Data

		entries := make([]dataSource.Entry, 0)
		for _, entr := range data {
			person := dataSource.NewPerson()
			for _, e := range entr {
				newEntry := dataSource.NewEntry(e.Key, e.Value)
				entries = append(entries, newEntry)
			}

			person.CountryID = country.ID
			person.Entries = entries

			if err := dataSource.SavePerson(&person); err != nil {
				tx.Rollback()

				return
			}
		}
	}

	tx.Commit()
}

func createCountries(list []string) (map[string]dataSource.Country, error) {
	countryList := make(map[string]dataSource.Country, 0)

	tx := dataSource.Transaction()
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

	tx.Commit()

	return countryList, nil
}
