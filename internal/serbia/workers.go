package serbia

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
			listing, err := common.GetListing(fmt.Sprintf("%s/nestali/lista/%d", baseUrl, page), "#missing-persons-main-container")

			if err != nil {
				producerStream <- nodeOrError{error: err}

				continue
			}

			items, err := common.QueryList(listing[0], ".osoba-main")

			if err != nil {
				producerStream <- nodeOrError{error: err}

				continue
			}

			if len(items) == 0 {
				stopFn()

				break
			}

			fmt.Printf("Serbia: Fetching page %d\n", page)

			for _, node := range items {
				producerStream <- nodeOrError{
					container: listing[0],
					node:      node,
					error:     nil,
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
			logger.Error("serbia", person.error.Error())
			fmt.Println(person.error.Error())

			return
		}

		*people = append(*people, person.person)
	}
}
