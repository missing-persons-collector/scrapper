package main

import (
	"fmt"
	"log"
	"missingPersons/dataSource"
	"missingPersons/logger"
	"sync"
)

func runDb() {
	if err := dataSource.NewDataSource("database", "postgres", "database", "database"); err != nil {
		log.Fatalf("Cannot connect to postgres database: %s", err.Error())
	}
}

func runLoggers() {
	fmt.Println("Creating loggers...")
	if err := logger.BuildLoggers([]string{"croatia", "serbia"}); err != nil {
		log.Fatalf("Unable to build loggers: %s\n", err.Error())
	}
	fmt.Println("Loggers created!")
}

func runCreateImageDir() {
	fmt.Println("Creating image directory...")
	createImageDirIfNotExists()
	fmt.Println("Image directory created!")
}

func runCountries() map[string]dataSource.Country {
	fmt.Println("Creating countries if they do not exist...")
	countryMap, err := createCountries([]string{"croatia", "serbia"})
	if err != nil {
		log.Fatalf("Error occurred while trying to create/find countries: %s. Exiting...", err.Error())
	}

	fmt.Println("Countries created or fetched. Continuing...\n")

	return countryMap
}

func runScrappers(countryMap map[string]dataSource.Country) {
	executions := createExecutions(countryMap)

	wg := &sync.WaitGroup{}
	for countryName, exec := range executions {
		wg.Add(1)
		go func(exec func() error, wg *sync.WaitGroup, countryName string) {
			if err := exec(); err != nil {
				//logger.Info(countryName, fmt.Sprintf("Country %s caused an error: %s. Continuing the rest of the countries...", countryName, err.Error()))
				fmt.Printf("Country %s caused an error: %s. Continuing the rest of the countries...\n", countryName, err.Error())
			}

			wg.Done()
		}(exec, wg, countryName)
	}

	wg.Wait()
}
