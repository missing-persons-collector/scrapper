package croatia

import (
	"fmt"
	"missingPersons/common"
	"regexp"
	"strings"
	"time"
)

func StartScrapping() ([]common.RawPerson, error) {
	page := 1
	baseUrl := "https://nestali.gov.hr"

	fieldMap := common.BuildFieldMap()

	people := make([]common.RawPerson, 0)
	for {
		listing, err := common.GetListing(fmt.Sprintf("%s/nestale-osobe-403/403?&page=%d", baseUrl, page), ".nestali-list .osoba-img")

		if err != nil {
			return nil, err
		}

		if len(listing) == 0 {
			break
		}

		fmt.Println(fmt.Sprintf("Croatia: Collecting page %d...", page))

		for _, item := range listing {
			href := common.GetAttr("href", item.Attr)
			url := fmt.Sprintf("%s%s", baseUrl, href)

			dataProperties, err := common.GetListing(url, ".profile_details_right dl *")

			if err != nil {
				return nil, err
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
					return nil, err
				}

				if field != "" {
					person, err = updateRawPerson(field, value, person)

					if err != nil {
						return nil, err
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

			people = append(people, person)
		}

		page++
	}

	return people, nil
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
