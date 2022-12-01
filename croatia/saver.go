package croatia

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"missingPersons/dataSource"
	"missingPersons/types"
	"regexp"
	"strings"
	"time"
)

func SaveCountry(pages []types.CollectedPage, country dataSource.Country) (types.Information, error) {
	fmt.Println("Croatia: Saving to database...")
	db := dataSource.DB()

	info := types.Information{
		UpdatedCount: 0,
		CreatedCount: 0,
		DeletedCount: 0,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, page := range pages {
			data := page.Data

			for _, entr := range data {
				id, err := createPersonId(entr)

				if err != nil {
					return err
				}

				var person dataSource.Person
				if err := db.Where("custom_id = ?", id).First(&person).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				fields, err := buildFields(entr)
				if err != nil {
					return err
				}

				if person.ID == "" {
					person = dataSource.NewPerson(id)
				}

				person, err = createPersonFromFields(id, fields, person)

				if err != nil {
					return err
				}

				person.CountryID = country.ID

				if person.ID == "" {
					if err := tx.Create(&person).Error; err != nil {
						return err
					}

					info.CreatedCount++
				} else {
					if err := tx.Save(&person).Error; err != nil {
						return err
					}

					info.UpdatedCount++
				}
			}
		}

		return nil
	}); err != nil {
		return info, err
	}

	fmt.Println("Croatia: All records saved to database.")

	return info, nil
}

func buildFieldMap() map[string]string {
	fieldMap := make(map[string]string)
	fieldMap["Ime"] = "Name"
	fieldMap["Prezime"] = "LastName"
	fieldMap["Djevojačko prezime"] = "MaidenName"
	fieldMap["Spol"] = "Gender"
	fieldMap["Datum rođenja"] = "DOB"
	fieldMap["Mjesto rođenja"] = "POB"
	fieldMap["Državljanstvo"] = "Citizenship"
	fieldMap["Prebivalište"] = "PrimaryAddress"
	fieldMap["Boravište"] = "SecondaryAddress"
	fieldMap["Država"] = "Country"
	fieldMap["Visina"] = "Height"
	fieldMap["Kosa"] = "Hair"
	fieldMap["Boja očiju"] = "EyeColor"
	fieldMap["Datum nestanka"] = "DOD"
	fieldMap["Mjesto nestanka"] = "POD"
	fieldMap["Okolnosti nestanka"] = "Description"

	return fieldMap
}

func buildFields(entries []types.ReceiverData) (map[string]string, error) {
	fieldMap := buildFieldMap()
	fieldsToUpdate := make(map[string]string)
	for _, entry := range entries {
		key := entry.Key

		if key == "Visina" {
			fmt.Println(entry)
		}

		for k, v := range fieldMap {
			matched, err := regexp.MatchString(k, key)

			if err != nil {
				return nil, err
			}

			if matched {
				fieldsToUpdate[v] = entry.Value
			}
		}
	}

	return fieldsToUpdate, nil
}

func createPersonFromFields(id string, fields map[string]string, person dataSource.Person) (dataSource.Person, error) {
	for k, v := range fields {
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
	}

	return person, nil
}

func createPersonId(entries []types.ReceiverData) (string, error) {
	keys := []string{"Ime", "Prezime", "Datum", "Mjesto rođenja"}
	values := make([]string, 0)

	for _, k := range keys {
		for _, e := range entries {
			matched, err := regexp.MatchString(k, e.Key)

			if err != nil {
				return "", err
			}

			if matched {
				values = append(values, e.Value)
			}
		}
	}

	final := ""
	for _, v := range values {
		final += v
	}
	final = strings.ToLower(final)

	re := regexp.MustCompile(`\s+`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`\.`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`,`)
	final = re.ReplaceAllString(final, "")

	return final, nil
}
