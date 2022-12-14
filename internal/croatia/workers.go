package croatia

import (
	"fmt"
	"missingPersons/common"
	"missingPersons/logger"
	worker2 "missingPersons/worker"
)

func producerFactory(baseUrl string) func(producerStream chan<- nodeOrError, stopFn func()) {
	return func(producerStream chan<- nodeOrError, stopFn func()) {
		page := 1

		for {
			listing, err := common.GetListing(fmt.Sprintf("%s/nestale-osobe-403/403?&page=%d", baseUrl, page), ".nestali-list .osoba-img")

			if err != nil {
				producerStream <- nodeOrError{error: err}

				continue
			}

			if len(listing) == 0 {
				stopFn()

				break
			}

			fmt.Printf("Croatia: Fetching page %d\n", page)

			for _, node := range listing {
				producerStream <- nodeOrError{
					node:  node,
					error: nil,
				}
			}

			page++
		}
	}
}

func consumerFactory(baseUrl string, fieldMap map[string]string) func(val interface{}, consumerStream chan<- personOrError) {
	return func(val interface{}, consumerStream chan<- personOrError) {
		processPerson(baseUrl, fieldMap, val.(nodeOrError), consumerStream)
	}
}

func waitFactory(people *[]common.RawPerson) func(data worker2.DataOrError) {
	return func(data worker2.DataOrError) {
		person := data.(personOrError)
		if person.error != nil {
			logger.Error("croatia", person.error.Error())
			fmt.Println(person.error.Error())

			return
		}

		*people = append(*people, person.person)
	}
}
