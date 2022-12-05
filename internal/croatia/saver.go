package croatia

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"missingPersons/common"
	"missingPersons/dataSource"
	"missingPersons/download"
	"missingPersons/types"
	"regexp"
	"strings"
)

func SaveCountry(people []common.RawPerson, country dataSource.Country, imageSaver download.ImageSaver) (types.Information, error) {
	fmt.Println("Croatia: Saving to database...")
	db := dataSource.DB()

	info := types.Information{
		UpdatedCount: 0,
		CreatedCount: 0,
		DeletedCount: 0,
	}

	customIds := make([]string, 0)
	if err := db.Transaction(func(tx *gorm.DB) error {
		for i, person := range people {
			if i%100 == 0 {
				fmt.Printf("Processed %d entries...\n", i)
			}
			id, err := createPersonId(person)

			if err != nil {
				return err
			}

			var dbPerson dataSource.Person
			if err := db.Where("custom_id = ?", id).First(&dbPerson).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if dbPerson.ID == "" {
				dbPerson = common.PersonFromRawPerson(id, country.ID, person)
			}

			if err != nil {
				return err
			}

			if person.ImageURL != "" {
				err := imageSaver.Save(person.ImageURL, id)

				if err != nil {
					fmt.Println(fmt.Sprintf("Cannot download and save image: %s", err.Error()))
				} else {
					dbPerson.ImageID = dbPerson.CustomID
				}
			}

			if dbPerson.ID == "" {
				if err := tx.Create(&dbPerson).Error; err != nil {
					return err
				}

				info.CreatedCount++
			} else {
				if err := tx.Save(&dbPerson).Error; err != nil {
					return err
				}

				info.UpdatedCount++
			}

			customIds = append(customIds, id)
		}

		return nil
	}); err != nil {
		return info, err
	}

	fmt.Println("Checking difference between the database and scraped data...")
	if err := diff(customIds); err != nil {
		fmt.Println(fmt.Sprintf("Difference could not be done: %s", err.Error()))
	}

	fmt.Println("Croatia: All records saved to database.")

	return info, nil
}

func diff(ids []string) error {
	offset := 1
	if err := dataSource.DB().Transaction(func(tx *gorm.DB) error {
		for {
			people := make([]dataSource.Person, 0)
			if err := dataSource.DB().Limit(100).Offset(offset).Find(&people).Error; err != nil {
				return err
			}

			if len(people) == 0 {
				break
			}

			for _, p := range people {
				lid := 0
				rid := len(ids) - 1

				found := false
				for lid <= rid {
					if ids[lid] == p.CustomID {
						found = true

						break
					}

					if ids[rid] == p.CustomID {
						found = true

						break
					}

					lid++
					rid--
				}

				if !found {
					if err := dataSource.DB().Where("ID = ?", p.ID).Delete(&dataSource.Person{}).Error; err != nil {
						return err
					}
				}
			}

			offset += 100
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func createPersonId(person common.RawPerson) (string, error) {
	final := fmt.Sprintf("%s%s%s%s", person.Name, person.LastName, person.DOB, person.POB)
	final = strings.ToLower(final)

	re := regexp.MustCompile(`\s+`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`/`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`\.`)
	final = re.ReplaceAllString(final, "")

	re = regexp.MustCompile(`,`)
	final = re.ReplaceAllString(final, "")

	return final, nil
}
