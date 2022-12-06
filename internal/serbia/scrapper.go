package serbia

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

	worker := worker2.NewWorker[nodeOrError, personOrError](5)
	worker.Produce(producerFactory("https://www.nestalisrbija.rs"))
	worker.Consume(consumerFactory("https://www.nestalisrbija.rs", common.BuildSerbiaFieldMap()))
	worker.Wait(waitFactory(&people))

	return people, nil
}

func processPerson(baseUrl string, fieldMap map[string]string, node nodeOrError, personStreamCh chan personOrError) {
	if node.error != nil {
		personStreamCh <- personOrError{error: node.error}

		return
	}

	item := node.node

	link, err := common.Query(item, "a")

	if err != nil {
		personStreamCh <- personOrError{error: node.error}

		return
	}

	href := common.GetAttr("href", link.Attr)

	container, err := common.GetListing(href, "#missing-persons-single-container")

	if err != nil {
		personStreamCh <- personOrError{error: err}

		return
	}

	dataProperties, err := common.QueryList(container[0], ".list-group-item")

	if err != nil {
		personStreamCh <- personOrError{error: err}

		return
	}

	person := common.NewRawPerson()
	for i := 0; i < len(dataProperties); i++ {
		var key string
		v := dataProperties[i]

		matched, _ := regexp.MatchString("Lični podaci", v.FirstChild.Data)
		if matched {
			continue
		}

		keyNode, err := common.Query(v, "b")
		if err != nil {
			personStreamCh <- personOrError{error: err}

			return
		}

		if keyNode != nil {
			key, err = determineField(keyNode.FirstChild.Data, fieldMap)

			if err != nil {
				personStreamCh <- personOrError{error: err}

				return
			}

			valueNode, err := common.Query(v, "span")
			if err != nil {
				personStreamCh <- personOrError{error: err}

				return
			}

			if key != "" && valueNode != nil && valueNode.FirstChild != nil {
				person, err = updateRawPerson(key, valueNode.FirstChild.Data, person)

				if err != nil {
					personStreamCh <- personOrError{error: err}

					return
				}
			}
		}
	}

	imageNode, err := common.Query(container[0], ".missing-item-image img")

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot retrieve node: %s", err.Error()))
	}

	if imageNode != nil {
		person.ImageURL = common.GetAttr("src", imageNode.Attr)
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
		if strings.ToLower(v) == "ženski" {
			v = "F"
		}

		if strings.ToLower(v) == "muški" {
			v = "M"
		}

		person.Gender = v
	}

	if k == "DOB" {
		t1 := strings.TrimRight(v, ". godine")
		t2 := strings.TrimRight(t1, ".")
		t2 = strings.Trim(t2, " ")

		months := map[string]string{
			"januar":  "1.",
			"februar": "2.",
			"mart":    "3.",
			"april":   "4.",
			"maj":     "5.",
			"jun":     "6.",
			"jul":     "7.",
			"avgust":  "8.",
			"septemb": "9.",
			"oktob":   "10.",
			"novemb":  "11.",
			"decemb":  "12.",
		}

		for k, v := range months {
			matched, _ := regexp.MatchString(k, t2)

			if matched {
				t := regexp.MustCompile(fmt.Sprintf("%s.[A-Za-z]?", k))
				t2 = t.ReplaceAllString(t2, v)
			}
		}

		t := regexp.MustCompile(`\s+`)
		t2 = t.ReplaceAllString(t2, "")

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

	if k == "Weight" {
		person.Weight = v
	}

	if k == "DOD" {
		t1 := strings.TrimRight(v, ". godine")
		t2 := strings.TrimRight(t1, ".")
		t2 = strings.Trim(t2, " ")

		months := map[string]string{
			"januar":  "1.",
			"februar": "2.",
			"mart":    "3.",
			"april":   "4.",
			"maj":     "5.",
			"jun":     "6.",
			"jul":     "7.",
			"avgust":  "8.",
			"septemb": "9.",
			"oktob":   "10.",
			"novemb":  "11.",
			"decemb":  "12.",
		}

		for k, v := range months {
			matched, _ := regexp.MatchString(k, t2)

			if matched {
				t := regexp.MustCompile(fmt.Sprintf("%s.[A-Za-z]?", k))
				t2 = t.ReplaceAllString(t2, v)
			}
		}

		t := regexp.MustCompile(`\s+`)
		t2 = t.ReplaceAllString(t2, "")

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
