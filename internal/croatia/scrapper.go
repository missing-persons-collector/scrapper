package croatia

import (
	"fmt"
	"missingPersons/common"
	worker2 "missingPersons/worker"
	"regexp"
	"strings"
	"time"
)

func StartScrapping() ([]common.RawPerson, error) {
	people := make([]common.RawPerson, 0)

	worker := worker2.NewWorker[nodeOrError, personOrError](20)
	worker.Produce(producerFactory("https://nestali.gov.hr"))
	worker.Consume(consumerFactory("https://nestali.gov.hr", common.BuildFieldMap()))
	worker.Wait(waitFactory(&people))

	return people, nil
}

func processPerson(baseUrl string, fieldMap map[string]string, node nodeOrError, personStreamCh chan personOrError) {
	if node.error != nil {
		personStreamCh <- personOrError{error: node.error}

		return
	}

	item := node.node

	href := common.GetAttr("href", item.Attr)
	url := fmt.Sprintf("%s%s", baseUrl, href)

	dataProperties, err := common.GetListing(url, ".profile_details_right dl *")

	if err != nil {
		personStreamCh <- personOrError{error: err}

		return
	}

	person := common.NewRawPerson()
	for i := 0; i < len(dataProperties); i++ {
		var key, value string
		v := dataProperties[i]

		if v.Data == "dt" {
			key = v.FirstChild.Data
			v := dataProperties[i+1]
			value = v.FirstChild.Data
			i++
		}

		field, err := determineField(key, fieldMap)

		if err != nil {
			personStreamCh <- personOrError{error: err}

			return
		}

		if field != "" {
			person, err = updateRawPerson(field, value, person)

			if err != nil {
				personStreamCh <- personOrError{error: err}

				return
			}
		}
	}

	imageNode, err := common.Query(item, "img")

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot retrieve node: %s", err.Error()))
	}

	if imageNode != nil {
		src := common.GetAttr("src", imageNode.Attr)
		person.ImageURL = fmt.Sprintf("%s%s", baseUrl, src)
	}

	personStreamCh <- personOrError{person: person}
}

func determineField(key string, fieldMap map[string]string) (string, error) {
	for k, v := range fieldMap {
		matched, err := regexp.MatchString(k, key)

		if err != nil {
			return "", err
		}

		if matched {
			return v, nil
		}
	}

	return "", nil
}

func updateRawPerson(k string, v string, person common.RawPerson) (common.RawPerson, error) {
	if k == "Name" {
		person.Name = v
	}

	if k == "LastName" {
		person.LastName = v
	}

	if k == "MaidenName" {
		person.MaidenName = v
	}

	if k == "Gender" {
		if v == "Ž" {
			v = "F"
		}

		person.Gender = v
	}

	if k == "DOB" {
		t1 := strings.TrimRight(v, ". godine")
		t2 := strings.TrimRight(t1, ".")

		date, err := time.Parse("2.1.2006", strings.TrimRight(t2, "."))

		if err != nil {
			person.DOB = ""
		} else {
			person.DOB = date.String()
		}

	}

	if k == "POB" {
		person.POB = v
	}

	if k == "Citizenship" {
		person.Citizenship = v
	}

	if k == "Country" {
		person.Country = v
	}

	if k == "PrimaryAddress" {
		person.PrimaryAddress = v
	}

	if k == "SecondaryAddress" {
		person.SecondaryAddress = v
	}

	if k == "Height" {
		person.Height = v
	}

	if k == "Hair" {
		person.Hair = v
	}

	if k == "EyeColor" {
		person.EyeColor = v
	}

	if k == "DOD" {
		t1 := strings.TrimRight(v, ". godine")
		t2 := strings.TrimRight(t1, ".")
		date, err := time.Parse("2.1.2006", strings.TrimRight(t2, "."))

		if err != nil {
			person.DOD = ""
		} else {
			person.DOD = date.String()
		}
	}

	if k == "POD" {
		person.POD = v
	}

	if k == "Description" {
		person.Description = v
	}

	return person, nil
}
